package kvdb

import (
	id "github.com/multiverse-os/levelup/id"
)

type KV struct {
	Name  string
	Key   id.Id
	Value []byte
}

type kvdb []KV

func (self kvdb) Len() int { return len(self) }

func (self kvdb) Less(i, j int) bool {
	return self[i].Key.UInt32() < self[j].Key.UInt32()
}

func (self kvdb) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// Aliasing
func (self kvdb) Size() int { return self.Len() }
