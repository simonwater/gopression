package exec

import (
	"errors"
	"fmt"

	"github.com/simonwater/gopression/chk"
	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/functions/funmgr"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

const STACK_MAX = 256

type VM struct {
	stack       [STACK_MAX]values.Value
	stackTop    int
	chunkReader *chk.ChunkReader
	tracer      *util.Tracer
}

func NewVM(tracer *util.Tracer) *VM {
	return &VM{
		tracer: tracer,
	}
}

func (vm *VM) reset() {
	vm.stackTop = 0
	vm.chunkReader = nil
}

func (vm *VM) push(value values.Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) pop() values.Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) peek() values.Value {
	return vm.stack[vm.stackTop-1]
}

func (vm *VM) peekAt(distance int) values.Value {
	return vm.stack[vm.stackTop-1-distance]
}

func (vm *VM) Execute(chunk *chk.Chunk, env env.Environment) ([]*ExResult, error) {
	chunkReader := chk.NewChunkReader(chunk, vm.tracer)
	return vm.ExecuteWithReader(chunkReader, env)
}

func (vm *VM) ExecuteWithReader(chunkReader *chk.ChunkReader, env env.Environment) ([]*ExResult, error) {
	vm.reset()
	vm.chunkReader = chunkReader
	return vm.run(env)
}

// run 虚拟机主循环
func (vm *VM) run(env env.Environment) ([]*ExResult, error) {
	if vm.tracer != nil {
		vm.tracer.StartTimerWithMsg("运行虚拟机")
		defer vm.tracer.EndTimer("虚拟机运行结束")
	}

	results := []*ExResult{}
	expOrder := 0

	for {
		op, err := vm.readCode()
		if err != nil {
			return results, errors.New("读取操作码失败: " + err.Error())
		}

		switch op {
		case chk.OP_BEGIN:
			expOrder, err = vm.readInt()
			if err != nil {
				return results, errors.New("读取表达式顺序失败: " + err.Error())
			}

		case chk.OP_END:
			val := vm.pop()
			result := &ExResult{
				Value: &val,
				State: OK,
				Index: expOrder,
			}
			results = append(results, result)

		case chk.OP_CONSTANT:
			value, err := vm.readConstant()
			if err != nil {
				return results, err
			}
			vm.push(value)

		case chk.OP_POP:
			vm.pop()

		case chk.OP_NULL:
			val := values.NewNullValue()
			vm.push(val)

		case chk.OP_GET_GLOBAL:
			name, err := vm.readString()
			if err != nil {
				return results, err
			}
			val := env.GetOrDefault(name, values.NewNullValue())
			vm.push(val)

		case chk.OP_SET_GLOBAL:
			name, err := vm.readString()
			if err != nil {
				return results, err
			}
			env.Put(name, vm.peek())

		case chk.OP_GET_PROPERTY:
			name, err := vm.readString()
			if err != nil {
				return results, err
			}
			obj := vm.pop()
			if !obj.IsInstance() {
				return results, fmt.Errorf("只有实例对象有属性: %s", name)
			}
			instance := obj.AsInstance()
			prop, _ := instance.Get(name)
			vm.push(prop)

		case chk.OP_SET_PROPERTY:
			name, err := vm.readString()
			if err != nil {
				return results, err
			}
			obj := vm.pop()
			if !obj.IsInstance() {
				return results, fmt.Errorf("只有实例对象有属性: %s", name)
			}
			value := vm.peek()
			instance := obj.AsInstance()
			instance.Set(name, value)

		case chk.OP_ADD:
			if err := vm.binaryOp(values.PLUS); err != nil {
				return results, err
			}

		case chk.OP_SUBTRACT:
			if err := vm.binaryOp(values.MINUS); err != nil {
				return results, err
			}

		case chk.OP_MULTIPLY:
			if err := vm.binaryOp(values.STAR); err != nil {
				return results, err
			}

		case chk.OP_DIVIDE:
			if err := vm.binaryOp(values.SLASH); err != nil {
				return results, err
			}

		case chk.OP_MODE:
			if err := vm.binaryOp(values.PERCENT); err != nil {
				return results, err
			}

		case chk.OP_POWER:
			if err := vm.binaryOp(values.STARSTAR); err != nil {
				return results, err
			}

		case chk.OP_GREATER:
			if err := vm.binaryOp(values.GREATER); err != nil {
				return results, err
			}

		case chk.OP_GREATER_EQUAL:
			if err := vm.binaryOp(values.GREATER_EQUAL); err != nil {
				return results, err
			}

		case chk.OP_LESS:
			if err := vm.binaryOp(values.LESS); err != nil {
				return results, err
			}

		case chk.OP_LESS_EQUAL:
			if err := vm.binaryOp(values.LESS_EQUAL); err != nil {
				return results, err
			}

		case chk.OP_EQUAL_EQUAL:
			if err := vm.binaryOp(values.EQUAL_EQUAL); err != nil {
				return results, err
			}

		case chk.OP_BANG_EQUAL:
			if err := vm.binaryOp(values.BANG_EQUAL); err != nil {
				return results, err
			}

		case chk.OP_NOT:
			if err := vm.preUnaryOp(values.BANG); err != nil {
				return results, err
			}

		case chk.OP_NEGATE:
			if err := vm.preUnaryOp(values.MINUS); err != nil {
				return results, err
			}

		case chk.OP_CALL:
			name, err := vm.readString()
			if err != nil {
				return results, err
			}
			if err := vm.callFunction(name); err != nil {
				return results, err
			}

		case chk.OP_JUMP_IF_FALSE:
			offset, err := vm.readInt()
			if err != nil {
				return results, err
			}
			if !vm.peek().IsTruthy() {
				if err := vm.gotoOffset(offset); err != nil {
					return results, err
				}
			}

		case chk.OP_JUMP:
			offset, err := vm.readInt()
			if err != nil {
				return results, err
			}
			if err := vm.gotoOffset(offset); err != nil {
				return results, err
			}

		case chk.OP_RETURN:
			// 暂时不处理

		case chk.OP_EXIT:
			if vm.stackTop != 0 {
				return results, fmt.Errorf("虚拟机状态异常，栈顶位置为：%d", vm.stackTop)
			}
			return results, nil

		default:
			return results, fmt.Errorf("暂不支持的指令：%s", op)
		}
	}
}

func (vm *VM) callFunction(name string) error {
	funcObj := funmgr.GetFunctionManager().GetFunction(name)
	cnt := funcObj.Arity()
	args := make([]values.Value, cnt)
	for i := cnt - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}
	result, err := funcObj.Call(args)
	if err != nil {
		return err
	}
	vm.push(result)
	return nil
}

func (vm *VM) binaryOp(tokenType values.TokenType) error {
	b := vm.pop()
	a := vm.pop()
	result, err := values.BinaryOperate(a, b, tokenType)
	if err != nil {
		return err
	}
	vm.push(result)
	return nil
}

func (vm *VM) preUnaryOp(tokenType values.TokenType) error {
	operand := vm.pop()
	result, err := values.PreUnaryOperate(operand, tokenType)
	if err != nil {
		return err
	}
	vm.push(result)
	return nil
}

func (vm *VM) readString() (string, error) {
	value, err := vm.readConstant()
	if err != nil {
		return "", err
	}
	return value.AsString(), nil
}

func (vm *VM) readConstant() (values.Value, error) {
	index, err := vm.readInt()
	if err != nil {
		return values.NewNullValue(), err
	}
	return vm.chunkReader.ReadConst(index)
}

func (vm *VM) readCode() (chk.OpCode, error) {
	code, err := vm.readByte()
	if err != nil {
		return chk.OP_EXIT, err
	}
	return chk.OpCodeFromValue(code)
}

func (vm *VM) readByte() (byte, error) {
	b, err := vm.chunkReader.ReadByte()
	return b, err
}

func (vm *VM) readInt() (int, error) {
	i, err := vm.chunkReader.ReadInt()
	return int(i), err
}

func (vm *VM) gotoOffset(offset int) error {
	curPos := vm.chunkReader.Position()
	return vm.chunkReader.NewPosition(curPos + offset)
}
