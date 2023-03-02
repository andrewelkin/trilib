package utils

import (
	"sync/atomic"
)

// emulated Atomic bool
type AtomBool struct {
	flag int32
}

func (b *AtomBool) Set() {
	atomic.StoreInt32(&(b.flag), 1)
}

func (b *AtomBool) Reset() {
	atomic.StoreInt32(&(b.flag), 0)
}

func (b *AtomBool) Get() bool {
	if atomic.LoadInt32(&(b.flag)) != 0 {
		return true
	}
	return false
}

func (b *AtomBool) GetAndSwap(newVal bool) bool {
	n := int32(0)
	if newVal {
		n = int32(1)
	}
	if atomic.SwapInt32(&(b.flag), n) != 0 {
		return true
	}
	return false
}
