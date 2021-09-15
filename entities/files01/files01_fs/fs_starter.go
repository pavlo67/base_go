package files01_fs

import (
	"fmt"
	"log"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/files01"
)

func Starter() starter.Operator {
	return &filesFSStarter{}
}

var l logger.Operator
var _ starter.Operator = &filesFSStarter{}

type filesFSStarter struct {
	access       config.Access
	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey

	// pathInfix    string
}

func (ffs *filesFSStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ffs *filesFSStarter) Prepare(cfg *config.Config, options common.Map) error {

	//ffs.basePath = strings.TrimSpace(options.StringDefault("base_path", ""))
	//if ffs.basePath == "" {
	//	return fmt.Errorf("no 'base_path' in options: %#v", options)
	//}

	configKey := strings.TrimSpace(options.StringDefault("config_key", "files_fs"))
	if configKey == "" {
		return fmt.Errorf("no 'config_key' in options (%#v)", options)
	}

	if err := cfg.Value(configKey, &ffs.access); err != nil {
		return errors.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	}

	log.Printf("%s --> %#v", configKey, ffs.access)

	ffs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(files01.InterfaceKey)))
	ffs.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(files01.InterfaceKeyCleaner)))

	return nil
}

func (ffs *filesFSStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	filesOp, filesCleanerOp, err := New(ffs.access.Path)
	if err != nil {
		return errors.Wrap(err, "can't init *filesFS{} as files.Operator")
	}

	if err = joinerOp.Join(filesOp, ffs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *filesFS{} as files.Operator with key '%s'", ffs.interfaceKey)
	}

	if err = joinerOp.Join(filesCleanerOp, ffs.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *filesFS{} as db.Cleaner with key '%s'", ffs.cleanerKey)
	}

	return nil
}
