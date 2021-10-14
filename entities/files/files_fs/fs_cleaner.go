package files_fs

import (
	"os"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
)

var _ db.Cleaner = &filesFS{}

const onClean = "on filesFS.Clean()"

func (filesOp *filesFS) Clean() error { // term *selectors.Term
	if err := os.RemoveAll(filesOp.basePath); err != nil {
		return errors.Wrapf(err, onClean+": removing %s", filesOp.basePath)
	}

	return nil
}
