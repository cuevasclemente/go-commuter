// commuters contains interfaces around data structures for
// which operations that commute with one another
// can be defined
package gocommuters

import "sync"

// commuter is a datastructure
// for which a commuting operation can be
// defined. Commuting operations are
// operations for which the order the
// operations are applied do not matter.
// If the semantics of the commuting operation
// require locking, it is recommended that
// you encode that logic into COp. I.E:
// func (c *MyCommuter) COp(s interface{}) {
//   c.Lock()
//   defer c.Unlock()
//   c.ReadCode(s.(string))
//   }
type Commuter interface {
	// COp is the commuting
	// operation for this datastructure
	COp(interface{})
	// Push queues a commutating operator
	Push(interface{})
	// Pop pops a datum for COp to be
	// performed on
	Pop() interface{}
	// EmptyQueue returns the queue of the
	// commutating operator. It returns an
	// array of the data to be operated on
	// and empties the queue of the commuter
	EmptyQueue() []interface{}
	// GetQueueLength returns the length
	// of the queue
	GetQueueLength() int
}

// Dequeue simply pops an element from
// the commuter's queue, and then
// runs COp on the element
func (c *Commuter) Dequeue() {
	c.COp(c.Pop())
}

// AggregateOp aggregates a commuting operation
// to the queue for the commuter
func (c *Commuter) AggregateOp(i interface{}) {
	c.Enqueue(i)
}

// CollapseQueue runs all of the operations
// in the queue
func (c *Commuter) CollapseQueue() {
	q := c.EmptyQueue()
	for _, i := range q {
		c.COp(q)
	}
}

// CCollapseQueue runs all of the operations
// in the queue, running `numWorkers` goroutines
// to empty out the queue in a threaded way.
// CCollapseQueue does not block or have
// any logic around waiting for completion, it
// tries to be as asyncrhronous as possible
func (c *Commuter) PCollapseQueue(numWorkers int) {
	l := c.GetQueueLength()
	// integer division in Go
	// takes the floor of the
	// quotient
	s := l / numWorkers
	q := c.Queue()
	for i := 0; i < s; i++ {
		go func(c *Commuter, is []interface{}) {
			for _, i := range is {
				c.COp(i)
			}
		}(q[i*s : (i+1)*s])
	}
}

// PCollapseQueue runs all of the operations
// in the queue, running `numWorkers` goroutines
// to empty out the queue in a threaded way.
// It is exactly the same as CCollapseQueue
// but it waits for all processes to end
func (c *Commuter) PCollapseQueue(numWorkers int) {
	l := c.GetQueueLength()
	// integer division in Go
	// takes the floor of the
	// quotient
	s := l / numWorkers
	q := c.Queue()
	var wg sync.WaitGroup
	for i := 0; i < s; i++ {
		wg.Add(1)
		go func(c *Commuter, is []interface{}) {
			defer wg.Done()
			for _, i := range is {
				c.COp(i)
			}
		}(q[i*s : (i+1)*s])
	}
	wg.Wait()
}

// CutQueue runs enough operations in the queue
// to get it to length below `desiredLength`.
// If the queue is already of or below
// that certain size, then CompressQueue
// is a no-op
func (c *Commuter) CompressQueue(desiredLength int) {
	for c.QueueLength() > desiredLength {
		c.Dequeue()
	}
}

// PCutQueue runs enough operations in the queue
// to get it to length below `desiredLength`.
// If the queue is already of or below
// that certain size, then CompressQueue
// is a no-op. PCompressQueue runs using
// `numWorkers` threads to dequeue. PCompressQueue
// blocks until the queue has been properly
// compressed
func (c *Commuter) CompressQueue(desiredLength int) {
	for c.QueueLength() > desiredLength {
		c.Dequeue()
	}
}

// CCutQueue runs enough operations in the queue
// to get it to length below `desiredLength`.
// If the queue is already of or below
// that certain size, then CCompressQueue
// is close to a no-op. CCompressQueue runs using
// `numWorkers` goroutines to dequeue.
func (c *Commuter) CCompressQueue(desiredLength int, numWorkers int) {
	numOps := c.QueueLength() - desiredLength
	var j *int
	*j = 0
	// somewhat questionable whether
	// or not this is valid, will fix later
	for i := 0; i < numWorkers; i++ {
		go func(j *int) {
			for *j < numOps {
				c.Dequeue()
				*j++
			}
		}(j)
	}
}
