package cache

import (
	"testing"
)

func TestNew(t *testing.T) {
	// with no limit
	bc := New[int](nil)
	got := bc.count
	if got != 0 {
		t.Errorf("New (no limit) count = %d; want 0;", got)
	}
	if bc.limit != nil {
		t.Errorf("New (no limit) limit = %d; want nil;", bc.limit)
	}
}

func TestNewWithLimit(t *testing.T) {
	limit := 512
	bc := New[int](&limit)
	got := bc.count
	if got != 0 {
		t.Errorf("New (with limit) count = %d; want 0;", got)
	}
	if *bc.limit != limit {
		t.Errorf("New (with limit) limit = %d; want %d;", bc.limit, limit)
	}
}

func TestPush(t *testing.T) {
	item := 1
	bc := New[int](nil)

	bc.Push(item)

	got := bc.head.item
	if got != item {
		t.Errorf("TestPush item = %d; want %d;", got, item)
	}

	got = bc.items[0]
	if got != item {
		t.Errorf("TestPush items[0] = %d; want %d;", got, item)
	}
}

func TestPop(t *testing.T) {
	item := 1
	bc := New[int](nil)
	bc.Push(item)

	got, err := bc.Pop()
	if err != nil {
		t.Errorf("TestPop err = %s;", err)
	}
	// returned item is correct
	if got != item {
		t.Errorf("TestPop item = %d; want %d;", got, item)
	}
	// head is set nil if last item popped
	if bc.head != nil {
		t.Errorf("TestPop head = %d; want nil;", bc.head.item)
	}
	// count is 0 if last item popped
	if bc.count != 0 {
		t.Errorf("TestPop count = %d; want 0;", bc.count)
	}
	// items is empty if last item popped
	if len(bc.items) != 0 {
		t.Errorf("TestPop items = %d; want 0;", len(bc.items))
	}
	// with empty stack
	_, err = bc.Pop()
	if err == nil {
		t.Errorf("TestPop count = %d; want 0;", bc.count)
	}
}

func TestPeak(t *testing.T) {
	item := 1
	bc := New[int](nil)
	bc.Push(item)
	got, err := bc.Peak()
	if err != nil {
		t.Errorf("TestPeak item = %d; want %d;", got, item)
	}
	if bc.head.item != item {
		t.Errorf("TestPeak head = %d; want %d;", bc.head.item, item)
	}
}

func TestCount(t *testing.T) {
	bc := New[int](nil)
	got := bc.Count()

	// with empty stack
	if got != 0 {
		t.Errorf("TestCount (empty) count = %d; want 0;", got)
	}
	// with n items in stack
	want := 3
	for i := 0; i < want; i++ {
		bc.Push(i + 1)
	}
	got = bc.Count()
	if got != want {
		t.Errorf("TestCount count = %d; want %d;", got, want)
	}
	// count matches items length
	if got != len(bc.items) {
		t.Errorf("TestCount count = %d; want %d;", got, len(bc.items))
	}
}

func TestString(t *testing.T) {
	bc := New[int](nil)
	bc.Push(1)
	bc.Push(2)
	bc.Push(3)
	got := bc.String()
	want := "[3 2 1]"
	if got != want {
		t.Errorf("TestString string = %s; want %s;", got, want)
	}
}

func TestLimit(t *testing.T) {
	limit := 10
	bc := New[int](&limit)
	for i := 1; i <= limit+10; i++ {
		bc.Push(i)
	}
	if bc.count != limit {
		t.Errorf("TestLimit limit = %d; want %d;", bc.count, limit)
	}
	// check head
	if bc.head.item != 20 {
		t.Errorf("TestLimit head = %d; want 20;", bc.head.item)
	}
	if bc.items[0] != 20 {
		t.Errorf("TestLimit items[0] = %d; want 20;", bc.items[0])
	}
	// check tail
	if bc.tail.item != 11 {
		t.Errorf("TestLimit tail = %d; want 11;", bc.tail.item)
	}
	if bc.items[len(bc.items)-1] != 11 {
		t.Errorf("TestLimit items[len(items)-1] = %d; want 11;", bc.items[len(bc.items)-1])
	}
}
