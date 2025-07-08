package values

import (
	"errors"
	"fmt"
	"math"
)

type ValuesHelper struct{}

func BinaryOperate(left, right Value, typ TokenType) (Value, error) {
	switch typ {
	case PLUS:
		if (!left.IsNumber() && !left.IsString()) || (!right.IsNumber() && !right.IsString()) {
			return NewNullValue(), errors.New("operands must be numbers or strings")
		}
		if left.IsString() || right.IsString() {
			return NewStringValue(left.String() + right.String()), nil
		} else {
			if left.IsDouble() || right.IsDouble() {
				return NewDoubleValue(left.AsDouble() + right.AsDouble()), nil
			} else {
				return NewIntValue(left.AsInteger() + right.AsInteger()), nil
			}
		}
	case MINUS:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		if left.IsDouble() || right.IsDouble() {
			return NewDoubleValue(left.AsDouble() - right.AsDouble()), nil
		} else {
			return NewIntValue(left.AsInteger() - right.AsInteger()), nil
		}
	case STAR:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		if left.IsDouble() || right.IsDouble() {
			return NewDoubleValue(left.AsDouble() * right.AsDouble()), nil
		} else {
			return NewIntValue(left.AsInteger() * right.AsInteger()), nil
		}
	case SLASH:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		if right.IsInteger() && right.AsInteger() == 0 {
			return NewNullValue(), errors.New("division by zero is not allowed")
		}
		if left.IsDouble() || right.IsDouble() {
			return NewDoubleValue(left.AsDouble() / right.AsDouble()), nil
		} else {
			return NewIntValue(left.AsInteger() / right.AsInteger()), nil
		}
	case PERCENT:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		if left.IsDouble() || right.IsDouble() {
			return NewDoubleValue(math.Mod(left.AsDouble(), right.AsDouble())), nil
		} else {
			return NewIntValue(left.AsInteger() % right.AsInteger()), nil
		}
	case STARSTAR:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		return NewDoubleValue(math.Pow(left.AsDouble(), right.AsDouble())), nil
	case GREATER:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		return NewBooleanValue(left.AsDouble() > right.AsDouble()), nil
	case GREATER_EQUAL:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		return NewBooleanValue(left.AsDouble() >= right.AsDouble()), nil
	case LESS:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		return NewBooleanValue(left.AsDouble() < right.AsDouble()), nil
	case LESS_EQUAL:
		if err := checkNumberOperands(left, right); err != nil {
			return NewNullValue(), err
		}
		return NewBooleanValue(left.AsDouble() <= right.AsDouble()), nil
	case BANG_EQUAL:
		return NewBooleanValue(!left.Equals(right)), nil
	case EQUAL_EQUAL:
		return NewBooleanValue(left.Equals(right)), nil
	default:
		return NewNullValue(), nil // Null value
	}
}

func PreUnaryOperate(operand Value, typ TokenType) (Value, error) {
	switch typ {
	case BANG:
		truthy := operand.IsTruthy()
		return NewBooleanValue(!truthy), nil
	case MINUS:
		if err := checkNumberOperand(operand); err != nil {
			return NewNullValue(), err // Null value
		}
		if operand.IsInteger() {
			return NewIntValue(-operand.AsInteger()), nil
		} else {
			return NewDoubleValue(-operand.AsDouble()), nil
		}
	default:
		return NewNullValue(), errors.New("暂不支持的类型") // Null value
	}
}

func checkNumberOperand(operand Value) error {
	if operand.IsNumber() {
		return nil
	}
	return errors.New("operand must be a number")
}

func checkNumberOperands(left, right Value) error {
	if left.IsNumber() && right.IsNumber() {
		return nil
	}
	return fmt.Errorf("operands must be numbers left: %v, right: %v", left, right)
}
