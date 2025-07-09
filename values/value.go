package values

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type Value struct {
	v  interface{}
	vt ValueType
}

func NewNullValue() Value {
	return Value{vt: Vt_Null}
}

func NewIntValue(i int32) Value {
	return Value{v: i, vt: Vt_Integer}
}

func NewDoubleValue(d float64) Value {
	return Value{v: d, vt: Vt_Double}
}

func NewStringValue(s string) Value {
	return Value{v: s, vt: Vt_String}
}

func NewBooleanValue(b bool) Value {
	return Value{v: b, vt: Vt_Boolean}
}

// Instance 类型请自行定义
func NewInstanceValue(inst Instance) Value {
	return Value{v: inst, vt: Vt_Instance}
}

func (val Value) GetValue() interface{} {
	return val.v
}

func (val Value) GetValueType() ValueType {
	return val.vt
}

// 反序列化
func GetFrom(buf *bytes.Buffer) (Value, error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return NewNullValue(), err
	}
	vt, ok := ValueOf(tag)
	if !ok {
		panic(fmt.Sprintf("未知类型: %T", vt))
	}
	switch vt {
	case Vt_Integer:
		var i int32
		if err := binary.Read(buf, binary.BigEndian, &i); err != nil {
			return NewNullValue(), err
		}
		return NewIntValue(i), nil
	case Vt_Double:
		var d float64
		if err := binary.Read(buf, binary.BigEndian, &d); err != nil {
			return NewNullValue(), err
		}
		return NewDoubleValue(d), nil
	case Vt_String:
		var slen int16
		if err := binary.Read(buf, binary.BigEndian, &slen); err != nil {
			return NewNullValue(), err
		}
		b := make([]byte, slen)
		if _, err := buf.Read(b); err != nil {
			return NewNullValue(), err
		}
		return NewStringValue(string(b)), nil
	default:
		return NewNullValue(), errors.New("暂不支持的类型")
	}
}

// 计算序列化字节长度
func (val Value) GetByteSize() (int16, error) {
	switch val.vt {
	case Vt_Integer:
		return 5, nil // 1 byte type + 4 bytes int32
	case Vt_Double:
		return 9, nil // 1 + 8
	case Vt_String:
		b := []byte(val.AsString())
		if len(b) > math.MaxInt16 {
			panic("字符串超出最大长度")
		}
		return int16(len(b) + 3), nil // 1 type + 2 len + n bytes
	default:
		panic("暂不支持的类型")
	}
}

// 序列化
func (val Value) WriteTo(buf *bytes.Buffer) error {
	switch val.vt {
	case Vt_Integer:
		buf.WriteByte(byte(val.vt))
		return binary.Write(buf, binary.BigEndian, int32(val.AsInteger()))
	case Vt_Double:
		buf.WriteByte(byte(val.vt))
		return binary.Write(buf, binary.BigEndian, val.AsDouble())
	case Vt_String:
		b := []byte(val.AsString())
		buf.WriteByte(byte(val.vt))
		if err := binary.Write(buf, binary.BigEndian, int16(len(b))); err != nil {
			return err
		}
		_, err := buf.Write(b)
		return err
	default:
		return errors.New("暂不支持的类型")
	}
}

func (val Value) IsBoolean() bool  { return val.vt == Vt_Boolean }
func (val Value) IsDouble() bool   { return val.vt == Vt_Double }
func (val Value) IsInteger() bool  { return val.vt == Vt_Integer }
func (val Value) IsNumber() bool   { return val.vt == Vt_Integer || val.vt == Vt_Double }
func (val Value) IsString() bool   { return val.vt == Vt_String }
func (val Value) IsNull() bool     { return val.vt == Vt_Null }
func (val Value) IsInstance() bool { return val.vt == Vt_Instance }

func (val Value) IsTruthy() bool {
	if val.IsNull() {
		return false
	} else if val.IsBoolean() {
		return val.AsBoolean()
	} else if val.IsString() {
		return len(val.AsString()) > 0
	}
	return true
}

func (val Value) AsBoolean() bool {
	if b, ok := val.v.(bool); ok {
		return b
	}
	panic(fmt.Sprintf("无法将 %T 转换为 bool", val.v))
}

func (val Value) AsInteger() int32 {
	switch v := val.v.(type) {
	case int:
		return int32(v)
	case int32:
		return v
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	}
	panic(fmt.Sprintf("无法将 %T 转换为 int32", val.v))
}

func (val Value) AsDouble() float64 {
	switch v := val.v.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	}
	panic(fmt.Sprintf("无法将 %T 转换为 float64", val.v))
}

func (val Value) AsString() string {
	if s, ok := val.v.(string); ok {
		return s
	}
	return ""
}

// Instance 类型请自行定义
func (val Value) AsInstance() Instance {
	return val.v.(Instance)
}

func (val Value) String() string {
	if val.v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", val.v)
}

func (val Value) Equals(other Value) bool {
	if val.vt != other.vt {
		return false
	}
	switch val.vt {
	case Vt_Null:
		return true
	case Vt_Boolean:
		return val.AsBoolean() == other.AsBoolean()
	case Vt_Integer:
		return val.AsInteger() == other.AsInteger()
	case Vt_Double:
		return val.AsDouble() == other.AsDouble()
	case Vt_String:
		return val.AsString() == other.AsString()
	default:
		return false
	}
}
