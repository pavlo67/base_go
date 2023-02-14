package convertor

import (
	"github.com/pavlo67/data/entities"
	"github.com/pavlo67/data/entities/files"
	"github.com/pavlo67/data/entities/records"
)

// Operator is an abstraction of external entity that can be converted into different our ones.
type Operator interface {
	// Original returns the original entity "as is".
	Original() (interface{}, error)

	// Object returns an object the original entity is converted into.
	Object() (*entities.Data, error)

	// FlowItem returns an flow.Item the original entity is converted into.
	FlowItem() (*records.Record, error)

	// Files returns a filer list the original entity is converted into.
	Files() ([]files.File, error)
}
