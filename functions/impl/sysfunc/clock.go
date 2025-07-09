package sysfunc

import (
	"strconv"
	"time"

	"github.com/simonwater/gopression/functions"
	"github.com/simonwater/gopression/values"
)

type Clock struct {
	*functions.Function
}

func NewClock() *Clock {
	return &Clock{
		Function: functions.NewFunction("clock", "当前毫秒数", functions.SYSTEM_GROUP),
	}
}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(arguments []values.Value) (values.Value, error) {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return values.NewStringValue(strconv.FormatInt(t, 10)), nil
}
