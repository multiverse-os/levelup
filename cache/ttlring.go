package cachering

import (
	"container/ring"
	"sync"
	"time"
)

// GetFromRemoteCallback returns current value of key from remote source. It's the long way.
type GetFromRemoteCallback func(key string) interface{}

type TTLRing struct {
	store         map[string]*Item
	maxSize       int
	getFromRemote GetFromRemoteCallback
	mutex         sync.Mutex
	keyRing       *ring.Ring
	statsAgent    *StatsAgent
}

func New(getFromRemoteCallback GetFromRemoteCallback, maxSize int) *TTLRing {
	return &TTLRing{
		store:         make(map[string]*Item, maxSize),
		maxSize:       maxSize,
		getFromRemote: getFromRemoteCallback,
		keyRing:       ring.New(maxSize),
		statsAgent:    NewStatsAgent(),
	}
}

// This function must be called from write locked mutex
func (q *TTLRing) refreshFromRemote(key string, ttl time.Duration) *Item {

	content := q.getFromRemote(key)
	if content == nil {
		_ = q.statsAgent.CommitEvent(EVENT_NOT_FOUND)
		return nil
	}
	item := &Item{
		Content:      content,
		LifeDuration: ttl,
		BirthTime:    time.Now().UTC(),
	}

	// Cap check and put value
	_, keyExist := q.store[key]

	if !keyExist {
		if q.keyRing.Value != nil { // Ring is full
			// Remove element in LRU fashion
			delete(q.store, q.keyRing.Value.(string))
		}
		q.keyRing.Value = key
		q.keyRing = q.keyRing.Next()
	}

	q.store[key] = item

	return item
}

func (q *TTLRing) Get(key string, ttl time.Duration) interface{} {

	item, exists := q.store[key]

	if !exists {
		// wait for acquiring the write lock. All read operations will be completed after this line.
		q.mutex.Lock()

		item, exists = q.store[key]
		if !exists { // Still not exists. So It's my job to get value from remote
			_ = q.statsAgent.CommitEvent(EVENT_MISS)
			item = q.refreshFromRemote(key, ttl)
		}

		q.mutex.Unlock()
	} else {

		if item.IsExpired() {
			q.mutex.Lock() // wait for the write lock and all read operations are completed
			_ = q.statsAgent.CommitEvent(EVENT_EXPIRED)
			item, exists = q.store[key]
			if !exists { // Key is deleted until we take the lock. So return nil
				return nil
			}
			if item.IsExpired() { // If item still expired it's our my job to update the value from remote
				item = q.refreshFromRemote(key, ttl)
			}

			q.mutex.Unlock()
			// Otherwise item is not expired
		} else {
			_ = q.statsAgent.CommitEvent(EVENT_HIT)
		}
	}
	return item.Content
	//from here we have the unexpired item, so return it.

}

func (q *TTLRing) GetKeyDuration(key string) *time.Duration {

	item, exists := q.store[key]

	if !exists {
		return nil
	}

	return &item.LifeDuration
}

func (q *TTLRing) GetKeyBirthTime(key string) *time.Time {

	item, exists := q.store[key]

	if !exists {
		return nil
	}

	return &item.BirthTime
}
