package levelup

import "time"

type EventType int

const (
	None EventType = iota
	Close
	Shutdown
)

type Action int

const (
	CreateEvent Action = iota
	UpdateEvent
	DeleteEvent
	RollbackEvent
)

type EventScope int

const (
	RecordScope EventScope = iota
	CollectionScope
	DatabaseScope
)

type Event struct {
	CreatedAt time.Time

	Action Action
	Type   EventType

	Record     *Record
	Collection *Collection
	Database   *Database

	//Differences []Difference

	// Rollback
	// TODO: Later we can simplify by removing rollback history if required
	//RewindEvent *Log
}

//type Difference struct {
//	Offset   int
//	Length   int
//	NewValue []byte
//	OldValue []byte
//}

type Events struct {
	Load   func(database Database)
	Unload func(database Database)

	Writing   func(database Database) (event Event)
	Rewinding func(database Database) (event Event)
	Tick      func(now time.Time) (delay time.Duration, event Event)
}
