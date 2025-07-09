package chk

import (
	"fmt"
	"sync"
)

// OpCode 表示操作码类型
type OpCode byte

// 操作码常量定义
const (
	OP_CONSTANT      OpCode = iota // 0
	OP_NULL                        // 1
	OP_TRUE                        // 2
	OP_FALSE                       // 3
	OP_POP                         // 4
	OP_GET_LOCAL                   // 5
	OP_SET_LOCAL                   // 6
	OP_GET_GLOBAL                  // 7
	OP_DEFINE_GLOBAL               // 8
	OP_SET_GLOBAL                  // 9
	OP_GET_PROPERTY                // 10
	OP_SET_PROPERTY                // 11
	OP_EQUAL_EQUAL                 // 12
	OP_BANG_EQUAL                  // 13
	OP_GREATER                     // 14
	OP_GREATER_EQUAL               // 15
	OP_LESS                        // 16
	OP_LESS_EQUAL                  // 17
	OP_ADD                         // 18
	OP_SUBTRACT                    // 19
	OP_MULTIPLY                    // 20
	OP_DIVIDE                      // 21
	OP_MODE                        // 22
	OP_POWER                       // 23
	OP_NOT                         // 24
	OP_NEGATE                      // 25
	OP_JUMP                        // 26
	OP_JUMP_IF_FALSE               // 27
	OP_CALL                        // 28
	OP_BEGIN                       // 29
	OP_END                         // 30
	OP_RETURN                      // 31
	OP_EXIT                        // 32
)

var (
	opCodeTitles = map[OpCode]string{
		OP_CONSTANT:      "OP_CONSTANT",
		OP_NULL:          "OP_NULL",
		OP_TRUE:          "OP_TRUE",
		OP_FALSE:         "OP_FALSE",
		OP_POP:           "OP_POP",
		OP_GET_LOCAL:     "OP_GET_LOCAL",
		OP_SET_LOCAL:     "OP_SET_LOCAL",
		OP_GET_GLOBAL:    "OP_GET_GLOBAL",
		OP_DEFINE_GLOBAL: "OP_DEFINE_GLOBAL",
		OP_SET_GLOBAL:    "OP_SET_GLOBAL",
		OP_GET_PROPERTY:  "OP_GET_PROPERTY",
		OP_SET_PROPERTY:  "OP_SET_PROPERTY",
		OP_EQUAL_EQUAL:   "OP_EQUAL_EQUAL",
		OP_BANG_EQUAL:    "OP_BANG_EQUAL",
		OP_GREATER:       "OP_GREATER",
		OP_GREATER_EQUAL: "OP_GREATER_EQUAL",
		OP_LESS:          "OP_LESS",
		OP_LESS_EQUAL:    "OP_LESS_EQUAL",
		OP_ADD:           "OP_ADD",
		OP_SUBTRACT:      "OP_SUBTRACT",
		OP_MULTIPLY:      "OP_MULTIPLY",
		OP_DIVIDE:        "OP_DIVIDE",
		OP_MODE:          "OP_MODE",
		OP_POWER:         "OP_POWER",
		OP_NOT:           "OP_NOT",
		OP_NEGATE:        "OP_NEGATE",
		OP_JUMP:          "OP_JUMP",
		OP_JUMP_IF_FALSE: "OP_JUMP_IF_FALSE",
		OP_CALL:          "OP_CALL",
		OP_BEGIN:         "OP_BEGIN",
		OP_END:           "OP_END",
		OP_RETURN:        "OP_RETURN",
		OP_EXIT:          "OP_EXIT",
	}

	valueToOpCode map[byte]OpCode
	opCodeMapOnce sync.Once
)

// 初始化操作码映射
func initOpCodeMap() {
	opCodeMapOnce.Do(func() {
		valueToOpCode = make(map[byte]OpCode)
		for op := OP_CONSTANT; op <= OP_EXIT; op++ {
			valueToOpCode[byte(op)] = op
		}
	})
}

// Value 返回操作码的字节值
func (op OpCode) Value() byte {
	return byte(op)
}

// Title 返回操作码的标题/助记符
func (op OpCode) Title() string {
	if title, ok := opCodeTitles[op]; ok {
		return title
	}
	return "unknown"
}

// String 实现Stringer接口
func (op OpCode) String() string {
	return fmt.Sprintf("%s(%d)", op.Title(), op)
}

// OpCodeForValue 根据字节值获取操作码
func OpCodeFromValue(value byte) (OpCode, error) {
	initOpCodeMap()
	if op, ok := valueToOpCode[value]; ok {
		return op, nil
	}
	return OP_EXIT, fmt.Errorf("未知的操作码: %d", value)
}
