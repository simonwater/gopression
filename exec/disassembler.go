package exec

import (
	"fmt"

	"github.com/simonwater/gopression/chk"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

type Disassembler struct {
	printer     func(msg string)
	chunkReader *chk.ChunkReader
}

// NewDisassembler 创建新的反汇编器
// printer: 自定义输出函数，如果为 nil 则使用默认 fmt.Println
func NewDisassembler(printer func(string)) *Disassembler {
	if printer == nil {
		printer = func(msg string) {
			fmt.Println(msg)
		}
	}
	return &Disassembler{printer: printer}
}

// Execute 执行反汇编过程
func (d *Disassembler) Execute(chunk *chk.Chunk) error {
	d.chunkReader = chk.NewChunkReader(chunk, util.NewTracer())
	var expOrder int32
	d.println("POSITION", "CODE", "PARAMETER", "ORDER")

	for {
		pos := fmt.Sprintf("%d", d.chunkReader.Position())
		op, err := d.readCode()
		if err != nil {
			return err
		}

		var param string
		switch op {
		case chk.OP_BEGIN:
			expOrder, err = d.chunkReader.ReadInt()
			if err != nil {
				return err
			}
			param = fmt.Sprintf("%d", expOrder)

		case chk.OP_CONSTANT:
			v, err := d.readConstant()
			if err != nil {
				return err
			}
			param = v.String()

		case chk.OP_GET_GLOBAL, chk.OP_SET_GLOBAL, chk.OP_GET_PROPERTY, chk.OP_SET_PROPERTY, chk.OP_CALL:
			param, err = d.readString()
			if err != nil {
				return err
			}

		case chk.OP_JUMP_IF_FALSE, chk.OP_JUMP:
			offset, err := d.chunkReader.ReadInt()
			if err != nil {
				return err
			}
			param = d.gotoOffset(int(offset))

		case chk.OP_EXIT:
			d.println(pos, op.Title(), "", fmt.Sprintf("%d", expOrder))
			return nil

		default:
			// 默认不处理参数
		}

		d.println(pos, op.Title(), param, fmt.Sprintf("%d", expOrder))
	}
}

func (d *Disassembler) readString() (string, error) {
	value, err := d.readConstant()
	return value.AsString(), err
}

func (d *Disassembler) readConstant() (values.Value, error) {
	index, err := d.chunkReader.ReadInt()
	if err != nil {
		return values.NewNullValue(), err
	}
	v, err := d.chunkReader.ReadConst(int(index))
	return *v, err
}

func (d *Disassembler) readCode() (chk.OpCode, error) {
	code, _ := d.chunkReader.ReadByte()
	return chk.OpCodeFromValue(code)
}

func (d *Disassembler) gotoOffset(offset int) string {
	curPos := d.chunkReader.Position()
	return fmt.Sprintf(":%d->to:%d", offset, curPos+offset)
}

func (d *Disassembler) println(pos, op, param, order string) {
	// 格式化输出行
	line := fmt.Sprintf("%-10s %-20s %-20s %s\n",
		truncate(pos, 10),
		truncate(op, 18),
		truncate(param, 18),
		order)

	d.printer(line)
}

// truncate 截断字符串到指定长度
func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}
