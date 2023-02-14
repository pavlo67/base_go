package persons_files

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db"
)

var _ db.Cleaner = &personsFiles{}

const onClean = "on personsFiles.Clean()"

func (personsOp *personsFiles) Clean() error {
	//if err := filelib.ClearDir(personsOp.path); err != nil {
	//	return errors.Wrap(err, onClean)
	//}
	return common.ErrNotImplemented
}
