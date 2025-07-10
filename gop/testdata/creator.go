package testdata

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/simonwater/gopression/env"
	"github.com/stretchr/testify/assert"
)

// 获取表达式列表
func GetExpressions(formulaBatches int) []string {
	lines := make([]string, 0, formulaBatches*5)

	fml := "A! = 1 + 2 * 3 - 6 - 1 + B! + C! * (D! - E! + 10 ** 2 / 5 - (12 + 8)) - F! * G! +  100 / 5 ** 2 ** 1"
	fml1 := "B! = C! + D! * 2 - 1"
	fml2 := "C! = D! * 2 + 1"
	fml3 := "D! = E! + F! * G!"
	fml4 := "G! = M! + N!"

	for i := 0; i < formulaBatches; i++ {
		idx := strconv.Itoa(i)
		lines = append(lines,
			strings.ReplaceAll(fml, "!", idx),
			strings.ReplaceAll(fml1, "!", idx),
			strings.ReplaceAll(fml2, "!", idx),
			strings.ReplaceAll(fml3, "!", idx),
			strings.ReplaceAll(fml4, "!", idx),
		)
	}

	return lines
}

// 创建环境并设置变量
func GetEnv(formulaBatches int) env.Environment {
	ev := env.NewDefaultEnvironment()

	for i := 0; i < formulaBatches; i++ {
		prefix := strconv.Itoa(i)
		ev.PutInt("E"+prefix, 2)
		ev.PutInt("F"+prefix, 3)
		ev.PutInt("M"+prefix, 4)
		ev.PutInt("N"+prefix, 5)
	}

	return ev
}

// 检查环境中的值
func CheckValues(t *testing.T, ev env.Environment, formulaBatches int) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		index := rand.Intn(formulaBatches)
		prefix := strconv.Itoa(index)

		a := ev.Get("A" + prefix)
		assert.Equal(t, 1686.0, a.GetValue(), "变量 A%s 值不匹配", prefix)

		b := ev.Get("B" + prefix)
		assert.Equal(t, 116, b.GetValue(), "变量 B%s 值不匹配", prefix)

		c := ev.Get("C" + prefix)
		assert.Equal(t, 59, c.GetValue(), "变量 C%s 值不匹配", prefix)

		d := ev.Get("D" + prefix)
		assert.Equal(t, 29, d.GetValue(), "变量 D%s 值不匹配", prefix)

		g := ev.Get("G" + prefix)
		assert.Equal(t, 9, g.GetValue(), "变量 G%s 值不匹配", prefix)
	}
}
