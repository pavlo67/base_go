package files_fs

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/files"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &catalogueFilesStarter{}
}

var l logger.Operator
var _ starter.Operator = &catalogueFilesStarter{}

type catalogueFilesStarter struct {
	fileKey      joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey

	// pathInfix    string
}

func (ffs *catalogueFilesStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ffs *catalogueFilesStarter) Prepare(cfg *config.Config, options common.Map) error {

	ffs.fileKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	ffs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(files.InterfaceKey)))
	ffs.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(files.InterfaceKeyCleaner)))

	return nil
}

func (ffs *catalogueFilesStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	filesOp, _ := joinerOp.Interface(ffs.fileKey).(files.Operator)
	if filesOp == nil {
		return fmt.Errorf("no files.Operator with key %s", ffs.fileKey)
	}

	catalogueOp, catalogueCleanerOp, err := New(filesOp)
	if err != nil {
		return errors.Wrap(err, "can't init *catalogueFiles{} as catalogue.Operator")
	}

	if err = joinerOp.Join(catalogueOp, ffs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *catalogueFiles{} as catalogue.Operator with key '%s'", ffs.interfaceKey)
	}

	if err = joinerOp.Join(catalogueCleanerOp, ffs.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *catalogueFiles{} as db.Cleaner with key '%s'", ffs.cleanerKey)
	}

	return nil
}
