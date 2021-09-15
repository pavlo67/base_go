package entities

import (
	"time"

	"github.com/pavlo67/data/elements/ns"
)

type Stat struct {
	NSS       ns.NSS
	ChldCount int64
	TotalSize int64
	CreatedAt time.Time
}
