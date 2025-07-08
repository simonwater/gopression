package values

import (
	"math"
	"testing"
)

func TestBinaryOperate_Plus(t *testing.T) {
	// int + int
	v, err := BinaryOperate(NewIntValue(2), NewIntValue(3), PLUS)
	if err != nil || !v.Equals(NewIntValue(5)) {
		t.Errorf("int+int failed: %v, %v", v, err)
	}
	// double + int
	v, err = BinaryOperate(NewDoubleValue(2.5), NewIntValue(3), PLUS)
	if err != nil || math.Abs(v.AsDouble()-5.5) > 1e-8 {
		t.Errorf("double+int failed: %v, %v", v, err)
	}
	// string + int
	v, err = BinaryOperate(NewStringValue("a"), NewIntValue(1), PLUS)
	if err != nil || v.String() != "a1" {
		t.Errorf("string+int failed: %v, %v", v, err)
	}
	// int + string
	v, err = BinaryOperate(NewIntValue(1), NewStringValue("b"), PLUS)
	if err != nil || v.String() != "1b" {
		t.Errorf("int+string failed: %v, %v", v, err)
	}
	// error
	v, err = BinaryOperate(NewBooleanValue(true), NewIntValue(1), PLUS)
	if err == nil {
		t.Error("expected error for bool+int")
	}
}

func TestBinaryOperate_MinusStarSlashPercent(t *testing.T) {
	// MINUS
	v, err := BinaryOperate(NewIntValue(5), NewIntValue(2), MINUS)
	if err != nil || !v.Equals(NewIntValue(3)) {
		t.Errorf("int-int failed: %v, %v", v, err)
	}
	v, err = BinaryOperate(NewDoubleValue(5.5), NewIntValue(2), MINUS)
	if err != nil || math.Abs(v.AsDouble()-3.5) > 1e-8 {
		t.Errorf("double-int minus failed: %v, %v", v, err)
	}
	// STAR
	v, err = BinaryOperate(NewIntValue(2), NewIntValue(3), STAR)
	if err != nil || !v.Equals(NewIntValue(6)) {
		t.Errorf("int*int failed: %v, %v", v, err)
	}
	// SLASH
	v, err = BinaryOperate(NewIntValue(6), NewIntValue(2), SLASH)
	if err != nil || !v.Equals(NewIntValue(3)) {
		t.Errorf("int/int failed: %v, %v", v, err)
	}
	v, err = BinaryOperate(NewDoubleValue(6.0), NewIntValue(4), SLASH)
	if err != nil || math.Abs(v.AsDouble()-1.5) > 1e-8 {
		t.Errorf("double/int slash failed: %v, %v", v, err)
	}
	// SLASH by zero
	_, err = BinaryOperate(NewIntValue(1), NewIntValue(0), SLASH)
	if err == nil {
		t.Error("expected division by zero error")
	}
	// PERCENT
	v, err = BinaryOperate(NewIntValue(7), NewIntValue(3), PERCENT)
	if err != nil || !v.Equals(NewIntValue(1)) {
		t.Errorf("int%%int failed: %v, %v", v, err)
	}
	v, err = BinaryOperate(NewDoubleValue(7.5), NewDoubleValue(2.0), PERCENT)
	if err != nil || math.Abs(v.AsDouble()-1.5) > 1e-8 {
		t.Errorf("double%%double failed: %v, %v", v, err)
	}
}

func TestBinaryOperate_StarStar(t *testing.T) {
	v, err := BinaryOperate(NewIntValue(2), NewIntValue(3), STARSTAR)
	if err != nil || math.Abs(v.AsDouble()-8.0) > 1e-8 {
		t.Errorf("2**3 failed: %v, %v", v, err)
	}
}

func TestBinaryOperate_Compare(t *testing.T) {
	// GREATER
	v, err := BinaryOperate(NewIntValue(3), NewIntValue(2), GREATER)
	if err != nil || !v.Equals(NewBooleanValue(true)) {
		t.Errorf("3>2 failed: %v, %v", v, err)
	}
	// LESS_EQUAL
	v, err = BinaryOperate(NewIntValue(2), NewIntValue(2), LESS_EQUAL)
	if err != nil || !v.Equals(NewBooleanValue(true)) {
		t.Errorf("2<=2 failed: %v, %v", v, err)
	}
}

func TestBinaryOperate_Equal(t *testing.T) {
	v, err := BinaryOperate(NewIntValue(2), NewIntValue(2), EQUAL_EQUAL)
	if err != nil || !v.Equals(NewBooleanValue(true)) {
		t.Errorf("2==2 failed: %v, %v", v, err)
	}
	v, err = BinaryOperate(NewIntValue(2), NewIntValue(3), BANG_EQUAL)
	if err != nil || !v.Equals(NewBooleanValue(true)) {
		t.Errorf("2!=3 failed: %v, %v", v, err)
	}
}

func TestPreUnaryOperate(t *testing.T) {
	// BANG
	v, err := PreUnaryOperate(NewBooleanValue(true), BANG)
	if err != nil || !v.Equals(NewBooleanValue(false)) {
		t.Errorf("!true failed: %v, %v", v, err)
	}
	// MINUS int
	v, err = PreUnaryOperate(NewIntValue(5), MINUS)
	if err != nil || !v.Equals(NewIntValue(-5)) {
		t.Errorf("-5 failed: %v, %v", v, err)
	}
	// MINUS double
	v, err = PreUnaryOperate(NewDoubleValue(2.5), MINUS)
	if err != nil || math.Abs(v.AsDouble()+2.5) > 1e-8 {
		t.Errorf("-2.5 failed: %v, %v", v, err)
	}
	// MINUS error
	_, err = PreUnaryOperate(NewStringValue("abc"), MINUS)
	if err == nil {
		t.Error("expected error for -string")
	}
}

func TestCheckNumberOperand(t *testing.T) {
	err := checkNumberOperand(NewIntValue(1))
	if err != nil {
		t.Errorf("int should be number")
	}
	err = checkNumberOperand(NewStringValue("a"))
	if err == nil {
		t.Error("string should not be number")
	}
}

func TestCheckNumberOperands(t *testing.T) {
	err := checkNumberOperands(NewIntValue(1), NewDoubleValue(2.0))
	if err != nil {
		t.Errorf("int and double should be numbers")
	}
	err = checkNumberOperands(NewIntValue(1), NewStringValue("a"))
	if err == nil {
		t.Error("int and string should not be numbers")
	}
}
