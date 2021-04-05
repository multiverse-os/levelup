package threadsafe

import "sync/atomic"

type Boolean uint32

func (self Boolean) Boolean() bool {
	return atomic.LoadUint32((*uint32)(&self))&1 == 1
}

func Switch() Boolean {
	return *new(Boolean)
}

func (self Boolean) On() {
	atomic.StoreUint32((*uint32)(&self), 1)
}

func (self Boolean) Off() {
	atomic.StoreUint32((*uint32)(&self), 0)
}

func (self Boolean) IsOn() bool  { return self.Boolean() }
func (self Boolean) IsOff() bool { return !self.Boolean() }

func (self Boolean) Flip(to bool) {
	atomic.StoreUint32((*uint32)(&self), booleanToUint32(to))
}

func (self Boolean) Toggle() bool {
	return atomic.AddInt32((int32)(&self), 1)&1 == 0
}

// SetToIf sets the Boolean to new only if the Boolean matches the old.
// Returns whether the set was done.
func (self Boolean) SetToIf(old, new bool) (set bool) {
	var o, n int32
	if old {
		o = 1
	}
	if new {
		n = 1
	}
	return atomic.CompareAndSwapInt32((*int32)(ab), o, n)
}

func booleanToUint32(boolValue bool) uint32 {
	return uint32(boolValue)
}

func uint32ToBoolean(val uint32) bool {
	if val > 0 {
		return true
	}
	return false
}
