package numfunc_test

import (
	"math"
	"testing"

	"github.com/simonwater/gopression/gop"
	"github.com/stretchr/testify/assert"
)

func TestAbsFunction(t *testing.T) {
	runner := gop.NewGopRunner()
	assert.Equal(t, 1, runner.Execute("abs(-1)"))
	assert.Equal(t, 1, runner.Execute("abs(1)"))
	assert.Equal(t, 1+2*3+int(math.Abs(float64(1-2*3))), runner.Execute("1 + 2 * 3 + abs(1 - 2 * 3)"))
	assert.Equal(t, 1+2*3+int(math.Abs(float64(1+2*3))), runner.Execute("1 + 2 * 3 + abs(1 + 2 * 3)"))
}
