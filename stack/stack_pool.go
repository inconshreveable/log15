// +build go1.3

package stack

import (
	"sync"
)

var pcStackPool = sync.Pool{
	New: func() interface{} { return make([]uintptr, 1000) },
}

// Callers returns a Trace for the current goroutine with element 0
// identifying the calling function.
func Callers() Trace {
	pcs := pcStackPool.Get().([]uintptr)
	pcs = pcs[:cap(pcs)]
	n := runtime.Callers(2, pcs)
	cs := make([]Call, n)
	for i, pc := range pcs[:n] {
		cs[i] = Call(pc)
	}
	pcStackPool.Put(pcs)
	return cs
}
