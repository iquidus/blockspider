package cache

// Node in chain (Doubly-linked)
type node[E any] struct {
	item E
	next *node[E]
	prev *node[E]
}

// The chain - a doubly-linked list with optional limit
type chain[E any] struct {
	head *node[E] // top of list
	tail *node[E] // bottom of list
	count int // number of nodes in list
	limit *int // max allowed nodes in list (nil = no max)
}

// Create a new chain and return a pointer to it
func newChain[E any](limit *int) *chain[E] {
	return &chain[E]{
		head: nil,
		tail: nil,
		count: 0,
		limit: limit,
	}
}

// Add node to top of chain
func (l *chain[E]) addToHead(n *node[E]) {
	if l.head == nil {
		l.tail = n
		l.head = n
	} else {
		n.next = l.head
		n.prev = nil
		n.next.prev = n
		l.head = n
	}
	l.count++
}

// Remove node from top of chain and return it
func (l *chain[E]) removeHead() *node[E] {
	if l.head == nil {
		return nil
	}
	// hold this for return
	oldHead := l.head

	// if theres a next node in chain, move head to it
	// otherwise set head as nil
	if l.head.next != nil {
		l.head = l.head.next
		l.head.prev = nil
	} else {
		l.head = nil
	}
	
	// remove from count
	l.count--
	return oldHead
}

// Remove node from bottom of chain
func (l *chain[E]) removeTail() *node[E] {
	n := l.tail
	if n == nil {
		return nil
	}
	l.tail = n.prev
	if l.tail != nil {
		l.tail.next = nil
	}
	l.count--
	return n
}