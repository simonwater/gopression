package values

import "fmt"

type ValueType byte

const (
	Vt_Integer  ValueType = 1
	Vt_Long     ValueType = 2
	Vt_Float    ValueType = 3
	Vt_Double   ValueType = 4
	Vt_String   ValueType = 5
	Vt_Boolean  ValueType = 6
	Vt_Instance ValueType = 7
	Vt_Null     ValueType = 8
)

var valueTypeMap = map[byte]ValueType{
	1: Vt_Integer,
	2: Vt_Long,
	3: Vt_Float,
	4: Vt_Double,
	5: Vt_String,
	6: Vt_Boolean,
	7: Vt_Instance,
	8: Vt_Null,
}

var valueTypeNames = map[ValueType]string{
	Vt_Integer:  "Integer",
	Vt_Long:     "Long",
	Vt_Float:    "Float",
	Vt_Double:   "Double",
	Vt_String:   "String",
	Vt_Boolean:  "Boolean",
	Vt_Instance: "Instance",
	Vt_Null:     "Null",
}

func (vt ValueType) Value() byte {
	return byte(vt)
}

func ValueOf(value byte) (ValueType, bool) {
	vt, ok := valueTypeMap[value]
	return vt, ok
}

// Title 返回枚举的标题字符串
func (vt ValueType) Title() string {
	if title, ok := valueTypeNames[vt]; ok {
		return title
	}
	panic(fmt.Sprintf("Unknown(%d)", vt))
}

// String 实现Stringer接口，调用Title
func (vt ValueType) String() string {
	return vt.Title()
}
