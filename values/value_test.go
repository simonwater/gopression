package values_test

import (
	"bytes"
	"testing"

	"github.com/simonwater/gopression/values"
)

func TestCreateNullValueByDefault(t *testing.T) {
	v := values.NewNullValue()
	if !v.IsNull() {
		t.Errorf("expected null value")
	}
	if v.GetValue() != nil {
		t.Errorf("expected value nil")
	}
	if v.GetValueType() != values.Vt_Null {
		t.Errorf("expected ValueType Null, got %v", v.GetValueType())
	}
}

func TestCreateIntegerAndDoubleValues(t *testing.T) {
	vi := values.NewIntValue(123)
	if !vi.IsInteger() {
		t.Errorf("expected integer")
	}
	if vi.AsInteger() != 123 {
		t.Errorf("expected 123, got %v", vi.AsInteger())
	}
	if vi.GetValueType() != values.Vt_Integer {
		t.Errorf("expected IntegerType")
	}
	vd := values.NewDoubleValue(3.14)
	if !vd.IsDouble() {
		t.Errorf("expected double")
	}
	if vd.AsDouble() != 3.14 {
		t.Errorf("expected 3.14, got %v", vd.AsDouble())
	}
	if vd.GetValueType() != values.Vt_Double {
		t.Errorf("expected DoubleType")
	}
}

func TestCreateStringAndBooleanValues(t *testing.T) {
	vs := values.NewStringValue("hello")
	if !vs.IsString() {
		t.Errorf("expected string")
	}
	if vs.AsString() != "hello" {
		t.Errorf("expected hello, got %v", vs.AsString())
	}
	if vs.GetValueType() != values.Vt_String {
		t.Errorf("expected StringType")
	}
	vb := values.NewBooleanValue(true)
	if !vb.IsBoolean() {
		t.Errorf("expected boolean")
	}
	if !vb.AsBoolean() {
		t.Errorf("expected true")
	}
	if vb.GetValueType() != values.Vt_Boolean {
		t.Errorf("expected BooleanType")
	}
}

func TestIsTruthy(t *testing.T) {
	if values.NewNullValue().IsTruthy() {
		t.Errorf("null should not be truthy")
	}
	if values.NewBooleanValue(false).IsTruthy() {
		t.Errorf("false should not be truthy")
	}
	if !values.NewBooleanValue(true).IsTruthy() {
		t.Errorf("true should be truthy")
	}
	if values.NewStringValue("").IsTruthy() {
		t.Errorf("empty string should not be truthy")
	}
	if !values.NewStringValue("abc").IsTruthy() {
		t.Errorf("non-empty string should be truthy")
	}
	if !values.NewIntValue(0).IsTruthy() {
		t.Errorf("0 should not be truthy")
	}
}

func TestInvalidAsXXX(t *testing.T) {
	v := values.NewStringValue("abc")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid asXXX")
		}
	}()
	_ = v.AsInteger()
	_ = v.AsDouble()
	_ = v.AsBoolean()
}

func TestEquals(t *testing.T) {
	if !values.NewIntValue(1).Equals(values.NewIntValue(1)) {
		t.Errorf("1 == 1 should be true")
	}
	if values.NewIntValue(1).Equals(values.NewIntValue(2)) {
		t.Errorf("1 == 2 should be false")
	}
	if !values.NewStringValue("a").Equals(values.NewStringValue("a")) {
		t.Errorf("'a' == 'a' should be true")
	}
	if values.NewStringValue("a").Equals(values.NewStringValue("b")) {
		t.Errorf("'a' == 'b' should be false")
	}
	if !values.NewNullValue().Equals(values.NewNullValue()) {
		t.Errorf("null == null should be true")
	}
	if values.NewBooleanValue(true).Equals(values.NewBooleanValue(false)) {
		t.Errorf("true == false should be false")
	}
}

func TestToString(t *testing.T) {
	if values.NewNullValue().String() != "null" {
		t.Errorf("expected 'null'")
	}
	if values.NewIntValue(123).String() != "123" {
		t.Errorf("expected '123'")
	}
	if values.NewStringValue("abc").String() != "abc" {
		t.Errorf("expected 'abc'")
	}
}

func TestWriteToAndGetFrom(t *testing.T) {
	vint := values.NewIntValue(42)
	vdouble := values.NewDoubleValue(3.14)
	vstr := values.NewStringValue("hi")
	buf := bytes.NewBuffer(nil)
	vint.WriteTo(buf)
	vdouble.WriteTo(buf)
	vstr.WriteTo(buf)
	rint, _ := values.GetFrom(buf)
	rdouble, _ := values.GetFrom(buf)
	rstr, _ := values.GetFrom(buf)
	if !rint.IsInteger() || rint.AsInteger() != 42 {
		t.Errorf("expected int 42")
	}
	if !rdouble.IsDouble() || rdouble.AsDouble() != 3.14 {
		t.Errorf("expected double 3.14")
	}
	if !rstr.IsString() || rstr.AsString() != "hi" {
		t.Errorf("expected string 'hi'")
	}
}

func TestGetFromUnsupportedType(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(99) // 非法tag
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unsupported type")
		}
	}()
	_, _ = values.GetFrom(buf)
}

func TestGetByteSizeOversizedString(t *testing.T) {
	v := values.NewStringValue(string(bytes.Repeat([]byte("a"), 40000)))
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for oversized string")
		}
	}()
	_, _ = v.GetByteSize()
}
