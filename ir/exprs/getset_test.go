package exprs_test

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/gop"
	"github.com/simonwater/gopression/values"
	"github.com/stretchr/testify/assert"
)

func TestGetSet_InstanceProperty(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	t1 := values.NewInstance()
	t1.Set("a", values.NewIntValue(1))
	t2 := values.NewInstance()
	t2.Set("b", values.NewIntValue(2))
	t2.Set("c", values.NewIntValue(3))
	environment.PutInstance("t1", t1)
	environment.PutInstance("t2", t2)

	runner := gop.NewGopRunner()
	lines := []string{
		"t1.x = t1.a + t2.b * t2.c + m",
		"m = t1.a + t2.b * t2.c",
	}
	result, _ := runner.ExecuteBatch(lines, environment)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 7, environment.Get("m").GetValue())
	t1Instance := environment.Get("t1").AsInstance()
	v, _ := t1Instance.Get("x")
	assert.Equal(t, 14, v.GetValue())
	assert.Equal(t, 14, result[0])
	assert.Equal(t, 7, result[1])
}
