package utils

type Queue[T any] struct {
	maxSize int
	Items   []T
}

func NewQueue[T any](maxSize int) *Queue[T] {
	return &Queue[T]{
		maxSize: maxSize,
		Items:   make([]T, 0, maxSize),
	}
}

func (q *Queue[T]) Enqueue(item T) {
	if len(q.Items) >= q.maxSize {
		q.Items = q.Items[1:]
	}
	q.Items = append(q.Items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.Items) == 0 {
		var zero T
		return zero, false
	}

	item := q.Items[0]
	q.Items = q.Items[1:]
	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if len(q.Items) == 0 {
		var zero T
		return zero, false
	}

	return q.Items[0], true
}

func (q *Queue[T]) Get() []T {
	return q.Items
}

func (q *Queue[T]) SetMaxSize(maxSize int) {
	q.maxSize = maxSize

	if len(q.Items) > maxSize {
		q.Items = q.Items[:maxSize]
	}
}

func (q *Queue[T]) Size() int {
	return len(q.Items)
}
