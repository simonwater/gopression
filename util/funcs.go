package util

import (
	"errors"
	"fmt"
)

func SafeExecute[T any](fn func() T) (result T, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	result = fn()
	return
}
