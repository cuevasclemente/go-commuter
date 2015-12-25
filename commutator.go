package gocommuter

// Commutator defines a way to commute two commuting operations.
// This is to say, Commutators define commutation rules for
// commuting operations on sets of data.
// For some discussion on this issue, let's say we defined
// a commuter:
// type MyCommuter struct {
//   baseValue float64
// }

// and a commutating operator:
// func (m *MyCommuter) COp(i interface{}) {
//   baseValue = baseValue * i.(float64)
// }

// a commutation rule would be useful if you had two values:

// v1 := 5.0
// v2 := -6.0

// and you were going to use COp(v1) and COp(v2), but you
// wanted to just call COp once. You would need a rule to know
// how different values would commute, and into what value they would
// commute.
// In this case, because multiplication is an operation that 'just'
// commutes, you would not need a commutator, the commutator for
// multiplication is
// is identity. On the other hand, if COp was instead:

// func (m *MyCommuter) COp(i interface{}) {
//   baseValue = baseValue - i.(float64)
// }

// The commutator would have to be

// myCommutator(i1 interface{}, i2 interface{}) {
//    i1 + i2
// }

// Note that:
// m.COp(v1); m.COp(v2) = (baseValue - 5.0) + 6.0  == baseValue + 1.0
// will yield the same value as
// m.COp(v2); m.COp(v1) = (baseValue + 6.0) - 5.0 == baseValue  + 1.0
// but:
// m.COp(v1-v2); = baseValue - 11
// and:
// m.COp(v2-v1); = baseValue + 11
//
// so despite the fact that our operation is
// defined in terms of subtraction, and it does
// indeed commute, (because the left-hand side stays
// constant), subtraction of arbitrary numbers does not commute.
// The commuter for "subtract by" is addition

// Note that:
// myCommutator(i1, i2) ==  1.0

// Essentially, Commutator ensures that
// m.COp(Commutator(v1, v2))
// results in the same operation as
// m.COp(v1); m.COp(v2)
type Commutator func(interface{}, interface{}) interface{}

type CommuterWithCommutator struct {
	Commuter
	Commutator
}

// CompressQueue compresses the queue of the commuter by
// acting with the commutator on pairs of elements
// popped from the queue and queuing the result to the
// Commuter. CompressQueue
// does this numOps times
// This means that compression can happen while
// the commuter is being accessed and queued.
func (c *CommuterWithCommutator) CompressQueue(numOps int) {
	for i := 0; i < numOps; i++ {
		i1 := c.Commuter.Pop()
		i2 := c.Commuter.Pop()
		c.Push(c.Commutator(i1, i2))
	}
}
