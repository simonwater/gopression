package numfunc

import (
	"errors"
	"math"

	"github.com/simonwater/gopression/functions"
	"github.com/simonwater/gopression/values"
)

type Abs struct {
	*functions.Function
}

func NewAbs() *Abs {
	return &Abs{
		Function: functions.NewFunction("abs", "绝对值", functions.NUMBER_GROUP),
	}
}

func (a *Abs) Arity() int {
	return 1
}

func (a *Abs) Call(arguments []values.Value) (values.Value, error) {
	if len(arguments) != 1 || !arguments[0].IsNumber() {
		panic(errors.New("参数不合法！"))
	}
	v := arguments[0]
	if v.IsDouble() {
		return values.NewDoubleValue(math.Abs(v.AsDouble())), nil
	}
	return values.NewIntValue(int32(math.Abs(float64(v.AsInteger())))), nil
}
