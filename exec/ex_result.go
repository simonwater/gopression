package exec

import "github.com/simonwater/gopression/values"

type ExResult struct {
	State ExState
	Value *values.Value
	Index int
	Error string
}

func NewExResult(value *values.Value, state ExState) *ExResult {
	return &ExResult{
		Value: value,
		State: state,
	}
}

func (r *ExResult) GetState() ExState {
	return r.State
}

func (r *ExResult) SetState(state ExState) {
	r.State = state
}

func (r *ExResult) GetResult() *values.Value {
	return r.Value
}

func (r *ExResult) SetResult(value *values.Value) {
	r.Value = value
}

func (r *ExResult) GetIndex() int {
	return r.Index
}

func (r *ExResult) SetIndex(index int) {
	r.Index = index
}

func (r *ExResult) GetError() string {
	return r.Error
}

func (r *ExResult) SetError(err string) {
	r.Error = err
}
