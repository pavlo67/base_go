package sources_stub

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/selectors"

	"github.com/pavlo67/workshop/components/sources"
)

var _ sources.Operator = &sourcesStub{}

// var _ db.Cleaner = &tagsSQLite{}

type sourcesStub struct {
	items []sources.Item
}

const onNew = "on sourcesStub.New(): "

func New(urls []string) (sources.Operator, db.Cleaner, error) {
	sourcesOp := sourcesStub{}
	for _, url := range urls {
		sourcesOp.items = append(sourcesOp.items, sources.Item{URL: url})
	}

	return &sourcesOp, nil, nil
}

func (sourcesOp *sourcesStub) Save(sources.Item, *crud.SaveOptions) (common.ID, error) {
	return "", common.ErrNotImplemented
}

func (sourcesOp *sourcesStub) Remove(common.ID, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

func (sourcesOp *sourcesStub) Read(common.ID, *crud.GetOptions) (*sources.Item, error) {
	return nil, common.ErrNotImplemented
}

func (sourcesOp *sourcesStub) List(*selectors.Term, *crud.GetOptions) ([]sources.Item, error) {
	return sourcesOp.items, nil
}

func (sourcesOp *sourcesStub) AddHistory(common.ID, crud.History, *crud.SaveOptions) (crud.History, error) {
	return nil, common.ErrNotImplemented
}
