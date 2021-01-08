package database

import (
	"time"

	"github.com/multiverse-os/levelup/database/model"
)

type LogType int

const (
	Create LogType = iota
	Update
	Delete
	Cache
)

type Log struct {
	CreatedAt time.Time

	Database   *Database
	Collection *Collection
	Record     *model.Record

	Type LogType
}
