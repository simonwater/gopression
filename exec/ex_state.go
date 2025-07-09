package exec

type ExState int

const (
	OK    ExState = 0
	ERROR ExState = 1
)

var exStateMap = map[int]ExState{
	0: OK,
	1: ERROR,
}

func (e ExState) ExStateValue() int {
	return int(e)
}

func ExStateFromValue(value int) (ExState, bool) {
	state, ok := exStateMap[value]
	return state, ok
}
