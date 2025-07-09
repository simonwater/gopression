package exec

import (
	"testing"

	"github.com/simonwater/gopression/env"
	"github.com/simonwater/gopression/parser"
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/visitors"
	"github.com/stretchr/testify/assert"
)

func execute(src string, e ...*env.DefaultEnvironment) any {
	var environment *env.DefaultEnvironment
	if len(e) > 0 {
		environment = e[0]
	} else {
		environment = env.NewDefaultEnvironment()
	}
	p := parser.NewParser(src)
	expr := p.Parse()

	tr := util.NewTracer()
	compiler := visitors.NewOpCodeCompiler(tr)
	compiler.BeginCompile()
	compiler.CompileExpr(expr, 0)
	chunk := compiler.EndCompile()

	vm := NewVM(tr)
	exResults, err := vm.Execute(chunk, environment)
	if err != nil {
		panic(err)
	}
	r := exResults[0].GetResult()
	if r != nil {
		if r.IsInteger() {
			return int(r.AsInteger())
		}
		return r.GetValue()
	}
	return nil
}

func TestVM_NumericalCalculations(t *testing.T) {
	assert.Equal(t, 7, execute("1 + 2 * 3"))
	assert.Equal(t, 0, execute("1 + 2 - 3"))
	assert.Equal(t, 9, execute("3 * (2 + 1)"))
	assert.Equal(t, 19.0, execute("1 + 2 * 3 ** 2 ** 1"))
	assert.Equal(t, 9.0, execute("3 * (2 + 1.0)"))
	assert.Equal(t, true, execute("3 * (2 + 1.0) > 7"))
}

func TestVM_VariableCalculations(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	environment.PutInt("a", 1)
	environment.PutInt("b", 2)
	environment.PutInt("c", 3)
	assert.Equal(t, 3.0, execute("a + b * c - 100 / 5 ** 2 ** 1", environment))
	assert.Equal(t, true, execute("a + b * c >= 6", environment))
	assert.Equal(t, 7, execute("x = y = a + b * c", environment))
	assert.Equal(t, 7, int(environment.Get("x").AsInteger()))
	assert.Equal(t, 7, int(environment.Get("y").AsInteger()))
}

func TestVM_StringConcatenation(t *testing.T) {
	assert.Equal(t, "hello world", execute(`"hello" + " world"`))
	assert.Equal(t, "a1", execute(`"a" + 1`))
	assert.Equal(t, "a1b", execute(`"a" + 1 + "b"`))
	assert.Equal(t, "a2", execute(`"a" + 1 * 2`))
}

func TestVM_Comparisons(t *testing.T) {
	assert.Equal(t, true, execute("1 < 2"))
	assert.Equal(t, true, execute("2 > 1"))
	assert.Equal(t, true, execute("1 <= 1"))
	assert.Equal(t, true, execute("2 >= 2"))
	assert.Equal(t, true, execute("1 == 1"))
	assert.Equal(t, true, execute("1 != 2"))
	assert.Equal(t, true, execute(`"a" == "a"`))
	assert.Equal(t, true, execute(`"a" != "b"`))
	assert.Equal(t, true, execute("1 < 2 && 3 > 2"))
}

func TestVM_UnaryOperations(t *testing.T) {
	assert.Equal(t, -1, execute("-1"))
	assert.Equal(t, -2.5, execute("-2.5"))
	assert.Equal(t, false, execute("!(1 == 1)"))
	assert.Equal(t, true, execute("!(1 == 2)"))
	assert.Equal(t, true, execute("!\"\""))
	assert.Equal(t, false, execute("!1"))
}

func TestVM_VariableAssignments(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	assert.Equal(t, 10, execute("x = 10", environment))
	assert.Equal(t, 10, int(environment.Get("x").AsInteger()))
	assert.Equal(t, 15, execute("y = x + 5", environment))
	assert.Equal(t, 15, int(environment.Get("y").AsInteger()))
	assert.Equal(t, 30, execute("z = y * 2", environment))
	assert.Equal(t, 30, int(environment.Get("z").AsInteger()))
}

func TestVM_NestedExpressions(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	environment.PutInt("a", 1)
	environment.PutInt("b", 2)
	assert.Equal(t, 7, execute("a + (b * 3)", environment))
	assert.Equal(t, 6, execute("(a + b) * 2", environment))
	assert.Equal(t, 5, execute("((a + b) * 2) - 1", environment))
}

func TestVM_ComplexExpressions(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	environment.PutInt("x", 10)
	environment.PutInt("y", 20)
	assert.Equal(t, 45, execute("x + y * 2 - 5", environment))
	assert.Equal(t, 15, execute("(x + y) / 2", environment))
	assert.Equal(t, 150, execute("x * (y - 5)", environment))
}

func TestVM_LogicalExpressions(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	environment.PutBool("a", true)
	environment.PutBool("b", false)
	assert.Equal(t, false, execute("a && b", environment))
	assert.Equal(t, true, execute("a || b", environment))
	assert.Equal(t, false, execute("!a", environment))
	assert.Equal(t, true, execute("!b", environment))
	assert.Equal(t, true, execute("a && (1 + 2 > 2)", environment))
}

func TestVM_MixedTypes(t *testing.T) {
	environment := env.NewDefaultEnvironment()
	environment.PutInt("a", 1)
	environment.PutString("b", "2")
	assert.Equal(t, "12", execute("a + b", environment))
	assert.Equal(t, "a2", execute(`"a" + b`, environment))
}
