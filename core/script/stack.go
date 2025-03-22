package script

type StackItem string

type Stack struct {
	Items []StackItem
}

func (s *Stack) Push(item StackItem) {
	s.Items = append(s.Items, item)
}

func (s *Stack) Pop() StackItem {
	item := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return item
}

func (s *Stack) Peek() StackItem {
	return s.Items[len(s.Items)-1]
}

func (s *Stack) Size() int {
	return len(s.Items)
}
