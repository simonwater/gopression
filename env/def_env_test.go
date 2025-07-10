package env

import (
	"testing"

	"github.com/simonwater/gopression/values"
	"github.com/stretchr/testify/assert"
)

func TestBaseEnvironment_Get_Put(t *testing.T) {
	env := NewDefaultEnvironment()
	env.Put("foo", values.NewStringValue("bar"))
	env.Put("baz", values.NewIntValue(42))
	env.Put("qux", values.NewBooleanValue(true))
	env.Put("quux", values.NewDoubleValue(3.14))
	env.PutString("testString", "Hello, World!")
	env.PutInt("testInt", 100)
	env.PutBool("testBool", true)
	env.PutDouble("testDouble", 3.14159)

	assert.Equal(t, "bar", env.Get("foo").GetValue())
	assert.Equal(t, 42, env.Get("baz").GetValue())
	assert.Equal(t, true, env.Get("qux").GetValue())
	assert.Equal(t, 3.14, env.Get("quux").GetValue())
	assert.Equal(t, "Hello, World!", env.Get("testString").AsString())
	assert.Equal(t, 100, int(env.Get("testInt").AsInteger()))
	assert.Equal(t, true, env.Get("testBool").AsBoolean())
	assert.Equal(t, 3.14159, env.Get("testDouble").AsDouble())
	assert.Equal(t, 8, env.Size(), "环境变量数量应为8")
}
