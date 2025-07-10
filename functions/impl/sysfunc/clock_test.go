package sysfunc_test

import (
	"fmt"
	"testing"

	"github.com/simonwater/gopression/gop"
)

func TestClockFunction(t *testing.T) {
	runner := gop.NewGopRunner()
	fmt.Println("Clock测试：")
	fmt.Println(runner.Execute("clock()"))
	fmt.Println(runner.Execute(`1 + 2 * 3 - 5 + " abc " + clock() + " " + 123`))
	fmt.Println("==========")
}
