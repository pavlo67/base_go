package files_fs

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/data/entities/files"
)

var _ files.Operator = &catalogueFiles{}

type catalogueFiles struct {
	filesOp files.Operator
}

const onNew = "on catalogueFiles.New(): "

func New(filesOp files.Operator) (files.Operator, db.Cleaner, error) {
	if filesOp == nil {
		return nil, nil, errors.New(onNew + ": no files.Operator")
	}

	catalogueOp := catalogueFiles{
		filesOp: filesOp,
	}

	return &catalogueOp, &catalogueOp, nil
}

const onSave = "on catalogueFiles.Save()"

func (catalogueOp *catalogueFiles) Save(path, newFilePattern string, data []byte, identity *auth.Identity) (string, error) {
	return catalogueOp.filesOp.Save(path, newFilePattern, data)
}

const onRead = "on catalogueFiles.Read()"

func (catalogueOp *catalogueFiles) Read(path string, identity *auth.Identity) ([]byte, error) {
	return catalogueOp.filesOp.Read(path)
}

const onRemove = "on catalogueFiles.Remove()"

func (catalogueOp *catalogueFiles) Remove(path string, identity *auth.Identity) error {
	return catalogueOp.filesOp.Remove(path)
}

const onList = "on catalogueFiles.Items()"

func (catalogueOp *catalogueFiles) List(path string, depth int, identity *auth.Identity) (files.Items, error) {
	return catalogueOp.filesOp.List(path, depth)
}

const onStat = "on catalogueFiles.Stat()"

func (catalogueOp *catalogueFiles) Stat(path string, depth int, identity *auth.Identity) (*files.Item, error) {
	return catalogueOp.filesOp.Stat(path, depth)

}
