package utils

// Simple Queue implementation
type Queue[T any] struct {
	maxSize int
	Items   []T
}

// Create a new Queue with a maximum size
func NewQueue[T any](maxSize int) *Queue[T] {
	return &Queue[T]{
		maxSize: maxSize,
		Items:   make([]T, 0, maxSize),
	}
}

// Add an item to the Queue
func (q *Queue[T]) Enqueue(item T) {
	if len(q.Items) >= q.maxSize {
		q.Items = q.Items[1:]
	}
	q.Items = append(q.Items, item)
}

// Remove an item from the Queue
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.Items) == 0 {
		var zero T
		return zero, false
	}

	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, true
}

// Get the first item in the Queue wtihout removing it
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.Items) == 0 {
		var zero T
		return zero, false
	}

	return q.Items[0], true
}

// Get all items in the Queue
func (q *Queue[T]) Get() []T {
	return q.Items
}

// Get the size of the Queue
func (q *Queue[T]) Size() int {
	return len(q.Items)
}

// Set the maximum size of the Queue
// If the new maximum size is smaller than the current size, the Queue will be truncated and items will potentially be lost
func (q *Queue[T]) SetMaxSize(maxSize int) {
	q.maxSize = maxSize

	if len(q.Items) > maxSize {
		q.Items = q.Items[:maxSize]
	}
}
