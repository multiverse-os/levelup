package cachering

import (
	"fmt"
)

type EventType string

const (
	EVENT_MISS EventType  = "miss"
	EVENT_HIT EventType = "hit"
	EVENT_EXPIRED EventType = "expired"
	EVENT_NOT_FOUND EventType = "not_found"
)


type StatsAgent struct {
	eventChan chan EventType
	counters map[EventType]int
}

func NewStatsAgent() *StatsAgent {
	a := &StatsAgent{
		eventChan: make(chan EventType, 256),
	}
	a.resetCounters()
	go func() {
		for {
			e := <- a.eventChan
			a.counters[e] += 1
		}
	}()
	return a
}

func (a *StatsAgent) resetCounters() {
	a.counters = map[EventType]int{
		EVENT_MISS: 0,
		EVENT_HIT: 0,
		EVENT_EXPIRED: 0,
		EVENT_NOT_FOUND: 0,
	}
}

func (a *StatsAgent) CommitEvent(eventType EventType) error {

	select {
	case a.eventChan <- eventType:
		return nil
	default:
		return fmt.Errorf("channel is full. Discarding event of %s", eventType)
	}

}

func (a *StatsAgent) GetStats() map[EventType]int {
	return a.counters
}