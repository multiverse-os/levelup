package history

import (
	"time"

	model "github.com/multiverse-os/levelup/model"
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

	Record *model.Record

	Type LogType
}
