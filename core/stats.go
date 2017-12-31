package core

import (
	"time"
)

type Stats struct {
	Rx uint64
	Tx uint64

	StartedAt  time.Time
	FinishedAt time.Time
}
