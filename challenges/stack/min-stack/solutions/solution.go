package solution

type entry struct {
	val    int
	minVal int
}

type MinStack struct {
	stack []entry
}

func NewMinStack() MinStack {
	return MinStack{}
}

func (s *MinStack) Push(val int) {
	minVal := val
	if len(s.stack) > 0 && s.stack[len(s.stack)-1].minVal < val {
		minVal = s.stack[len(s.stack)-1].minVal
	}
	s.stack = append(s.stack, entry{val, minVal})
}

func (s *MinStack) Pop() {
	s.stack = s.stack[:len(s.stack)-1]
}

func (s *MinStack) Top() int {
	return s.stack[len(s.stack)-1].val
}

func (s *MinStack) GetMin() int {
	return s.stack[len(s.stack)-1].minVal
}
