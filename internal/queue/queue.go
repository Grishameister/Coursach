package queue

import "github.com/enriquebris/goconcurrentqueue"

type Queue struct {
	queue *goconcurrentqueue.FixedFIFO
}

func NewQueue() *Queue {
	return &Queue{
		queue: goconcurrentqueue.NewFixedFIFO(100),
	}
}

func (q *Queue) Push(bytes []byte) error {
	return q.queue.Enqueue(bytes)
}

func (q *Queue) Pop() ([]byte, error) {
	elem, err := q.queue.DequeueOrWaitForNextElement()
	if err != nil {
		return nil, err
	}
	return elem.([]byte), nil
}

func (q *Queue) Size() int {
	return q.queue.GetLen()
}
