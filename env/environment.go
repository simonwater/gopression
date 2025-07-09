package env

import (
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

type Environment interface {
	BeforeExecute(vars []*util.Field) bool
	Get(id string) values.Value
	GetOrDefault(id string, defValue values.Value) values.Value
	Put(id string, value values.Value)
	Size() int
}
type BaseEnvironment struct {
	Environment
}

func NewBaseEnvironment(behavior Environment) *BaseEnvironment {
	return &BaseEnvironment{
		Environment: behavior,
	}
}

func (be *BaseEnvironment) PutInt(id string, value int32) {
	be.Put(id, values.NewIntValue(value))
}

// PutDouble 添加浮点型变量
func (be *BaseEnvironment) PutDouble(id string, value float64) {
	be.Put(id, values.NewDoubleValue(value))
}

// PutString 添加字符串变量
func (be *BaseEnvironment) PutString(id string, value string) {
	be.Put(id, values.NewStringValue(value))
}

// PutBool 添加布尔型变量
func (be *BaseEnvironment) PutBool(id string, value bool) {
	be.Put(id, values.NewBooleanValue(value))
}

// PutInstance 添加实例对象
func (be *BaseEnvironment) PutInstance(id string, obj *values.Instance) {
	be.Put(id, values.NewInstanceValue(*obj))
}
