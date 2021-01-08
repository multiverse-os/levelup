package kvdb

import (
	"path"
	"sort"
	"sync"

	id "github.com/multiverse-os/levelup/id"
)

//type Map struct {
//	mu sync.Mutex
//
//	read atomic.Value
//	dirty map[interface{}]*entry
//	misses int
//}
//
//type readOnly struct {
//	m       map[interface{}]*entry
//	amended bool
//}

type memory struct {
	*sync.Map
}

func NewMemoryDB() *memory {
	return &memory{
		new(sync.Map),
	}
}

func (self *memory) Put(key string, value []byte) {
	self.Store(key, KV{key, id.Hash(key), value})
}

func (s *memory) PutMany(m map[string]string) {
	for k, v := range m {
		s.Put(k, v)
	}
}

func (s *memory) Exists(key string) bool {
	if _, err := s.get(key); err != nil {
		return false
	}
	return true
}

func (s *memory) get(key string) (kvPair, error) {
	v, ok := s.Load(key)
	if !ok {
		return kvPair{}, ErrNotExist
	}
	return v.(kvPair), nil
}

func (s *memory) Get(key string, defaultValue ...string) (string, error) {
	kv, err := s.get(key)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", err
	}
	return kv.Value, nil
}

func (s *memory) getAllMatched(pattern string) (kvPairs, error) {
	kvs := make(kvPairs, 0)
	s.m.Range(func(_, value interface{}) bool {
		kv := value.(kvPair)
		if matched, _ := path.Match(pattern, kv.Key); matched {
			kvs = append(kvs, kv)
		}
		return true
	})
	if len(kvs) == 0 {
		return nil, ErrNoMatched
	}
	sort.Sort(kvs)
	return kvs, nil
}

func (s *memory) GetMany(pattern string) ([]string, error) {
	vs := make([]string, 0)
	kvs, err := s.getAllMatched(pattern)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		vs = append(vs, kv.Value)
	}
	sort.Strings(vs)
	return vs, nil
}

func (s *memory) Del(key string) { s.m.Delete(key) }

func (s *memory) Flush() {
	s.Range(func(key, _ interface{}) bool {
		s.Delete(key)
		return true
	})
}
