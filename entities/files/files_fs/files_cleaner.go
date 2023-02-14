package files_fs

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db"
)

var _ db.Cleaner = &catalogueFiles{}

const onClean = "on filesFS.Clean()"

func (filesOp *catalogueFiles) Clean() error {
	//if err := os.RemoveAll(filesOp.basePath); err != nil {
	//	return errors.Wrapf(err, onClean+": removing %s", filesOp.basePath)
	//}

	return common.ErrNotImplemented
}
