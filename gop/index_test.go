package gop

import (
	"fmt"
	"os"
	"testing"
)

func TestTest(t *testing.T) {
	printer := func(message string) {
		if message != "" {
			fmt.Fprintln(os.Stdout, message)
		}
	}
	printer("Hello, World!")
}
