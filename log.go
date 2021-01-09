package levelup

import (
	"time"
)

type ActionType int

const (
	CREATE ActionType = iota
	UPDATE
	DELETE
)

type ActionSubject int

const (
	DATABASE ActionSubject = iota
	COLLECTION
	RECORD
)

type ActionLog struct {
	CreatedAt time.Time

	Database   *Database
	Collection *Collection
	Record     *Record

	Action  ActionType
	Subject ActionSubject
	// TODO: Not only will this provide our desired rewinding functionality; it
	// will allow us to filter for creates on things like collections and records,
	// and rebuild the database on load in the read only cache form.
}
