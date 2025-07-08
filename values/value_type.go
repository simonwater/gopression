package values

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

func (vt ValueType) Value() byte {
	return byte(vt)
}

func ValueOf(value byte) (ValueType, bool) {
	vt, ok := valueTypeMap[value]
	return vt, ok
}
