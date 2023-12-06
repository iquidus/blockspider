package cache

import (
	"errors"
	"fmt"
)

type BlockStack[E any] struct {
	*chain[E]
}

// Create new stack and return a pointer to it
// Accepts a max item limit as *int, pass nil for no limit
func New[E any](limit *int) *BlockStack[E] {
	s := BlockStack[E]{
		chain: newChain[E](limit),
	}
	return &s
}

// Push item to head of stack.
// Removes tail if we are at limit
func (s *BlockStack[E]) Push(item E) {
	n := node[E]{
		item: item,
	}
	s.addToHead(&n)
	if s.limit != nil && s.count > *s.limit {
		s.removeTail()
	}
}

// Pop item from head of stack and return it, or return error.
func (s *BlockStack[E]) Pop() (E, error) {
	n := s.removeHead()
	if n == nil {
		var e E
		return e, errors.New("Stack is empty")
	}
	return n.item, nil
}

// Return item from head of stack, or return error.
func (s *BlockStack[E]) Peak() (E, error) {
	if s.head == nil {
		var e E
		return e, errors.New("Stack is empty")
	}
	return s.head.item, nil
}

// Return number of items in stack
func (s *BlockStack[E]) Count() int {
	return s.count
}

// Return items as a slice with head at index 0
func (s *BlockStack[E]) Items() []E {
	return s.items
}

// String function for fmt print etc.
func (s *BlockStack[E]) String() string {
	return fmt.Sprintf("%v", s.Items())
}
