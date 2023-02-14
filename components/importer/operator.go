package importer

import (
	"time"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/entities/records"
)

type DataSeries struct {
	URN     ns.URN
	At      time.Time
	Records []records.Record
	Notes   string
}

type Operator interface {
	// Prepare opens import session with selected data source
	// Init() error

	Get(urn ns.URN, options common.Map) (*DataSeries, error)
	IsNew(records.Record) (bool, error)
}
