package cachering

import (
	"time"
)

type Item struct {
	Content      interface{}
	LifeDuration time.Duration
	BirthTime    time.Time
}

func (item *Item) IsExpired() bool {
	deadline := item.BirthTime.Add(item.LifeDuration)
	return time.Now().UTC().After(deadline)
}