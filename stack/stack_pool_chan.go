// +build !go1.3

package stack

import (
	"runtime"
)

const (
	stackPoolSize = 64
)

type stackPool struct {
	c chan []uintptr
}

func newStackPool() *stackPool {
	return &stackPool{c: make(chan []uintptr, stackPoolSize)}
}

func (p *stackPool) Get() []uintptr {
	select {
	case st := <-p.c:
		return st
	default:
		return make([]uintptr, 1000)
	}
}

func (p *stackPool) Put(st []uintptr) {
	select {
	case p.c <- st:
	default:
	}
}

var pcStackPool = newStackPool()

// Callers returns a Trace for the current goroutine with element 0
// identifying the calling function.
func Callers() Trace {
	pcs := pcStackPool.Get()
	pcs = pcs[:cap(pcs)]
	n := runtime.Callers(2, pcs)
	cs := make([]Call, n)
	for i, pc := range pcs[:n] {
		cs[i] = Call(pc)
	}
	pcStackPool.Put(pcs)
	return cs
}
