package ir

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVarsQuery_Formula(t *testing.T) {
	varQuery := NewVarsQuery()
	result := varQuery.ExecuteSrc("x = y = a + b*(2 + (z = h * i)) - abs(sum(c, d - e/f**g))")
	assert.Equal(t, "x,y,z = a,b,c,d,e,f,g,h,i", result.String())
}

func TestVarsQuery_IfExpression(t *testing.T) {
	varQuery := NewVarsQuery()
	result := varQuery.ExecuteSrc("if(a == b + c, if (a > d, x = y = m + n, p = q = u + v), z = w * 2)")
	assert.Equal(t, "p,q,x,y,z = a,b,c,d,m,n,u,v,w", result.String())
}

func TestVarsQuery_FormulaWithInstance(t *testing.T) {
	varQuery := NewVarsQuery()
	result := varQuery.ExecuteSrc(
		"A.x = A.y = B.a + B.b*(2 + (A.z = C.D.h * C.D.i)) - abs(sum(B.c, B.d - C.D.e/C.D.f**C.D.g))",
	)
	assert.Equal(t, "A.x,A.y,A.z = B.a,B.b,B.c,B.d,C.D.e,C.D.f,C.D.g,C.D.h,C.D.i", result.String())
}

func TestVarsQuery_LargeNumberOfFormulas(t *testing.T) {
	t.Log("批量查询变量测试：")
	cnt := 10000
	t.Logf("公式总数：%d", cnt)
	lines := make([]string, 0, cnt)
	fml := "A! = 1 + 2 * 3 - 6 - 1 + B! + C! * (D! - E! + 10 ** 2 / 5 - (12 + 8)) - F! * G! +  100 / 5 ** 2 ** 1"
	for i := 1; i <= cnt; i++ {
		line := fml
		line = strings.ReplaceAll(line, "!", fmt.Sprintf("%d", i))
		lines = append(lines, line)
	}
	start := time.Now()
	varQuery := NewVarsQuery()
	var result *VariableSet
	for _, expr := range lines {
		result = varQuery.ExecuteSrc(expr)
	}
	t.Log(result)
	t.Logf("time: %dms", time.Since(start).Milliseconds())
	t.Log("==========")
}
