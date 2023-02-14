package persons_files

import (
	"fmt"

	"github.com/pavlo67/common/common/files"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons"
)

func Starter() starter.Operator {
	return &personsFilesStarter{}
}

const configKeyDefault = "persons_files"

type personsFilesStarter struct {
	filesKey            joiner.InterfaceKey
	interfaceKey        joiner.InterfaceKey
	interfaceCleanerKey joiner.InterfaceKey
}

var _ starter.Operator = &personsFilesStarter{}
var l logger.Operator

func (uks *personsFilesStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (uks *personsFilesStarter) Prepare(cfg *config.Config, options common.Map) error {

	uks.filesKey = joiner.InterfaceKey(options.StringDefault("files_key", string(files.InterfaceKey)))
	uks.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(persons.InterfaceKey)))
	uks.interfaceCleanerKey = joiner.InterfaceKey(options.StringDefault("interface_cleaner_key", string(persons.InterfaceCleanerKey)))

	return nil
}

func (uks *personsFilesStarter) Run(joinerOp joiner.Operator) error {

	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	filesOp, _ := joinerOp.Interface(uks.filesKey).(files.Operator)
	if filesOp == nil {
		return fmt.Errorf("no files.Operator with key %s", uks.filesKey)
	}

	personsOp, personsCleanerOp, err := New(filesOp)
	if err != nil {
		return errors.Wrapf(err, "can't personsFiles.New()")
	}

	if err = joinerOp.Join(personsOp, uks.interfaceKey); err != nil {
		return errors.Wrap(err, "can't join *personsFiles{} as persons.Operator interface")
	}

	if err = joinerOp.Join(personsCleanerOp, uks.interfaceCleanerKey); err != nil {
		return errors.Wrap(err, "can't join *personsFiles{} as db.Cleaner interface")
	}

	return nil
}
