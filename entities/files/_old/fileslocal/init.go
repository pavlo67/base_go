package fileslocal

import (
	"os"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"fmt"

	"github.com/pavlo67/punctum/things_old/fileinfo"
	"github.com/pavlo67/punctum/things_old/files"
)

func Starter() starter.Operator {
	return &fileslocalStarter{}
}

const InterfaceKey joiner.InterfaceKey = "fileslocal"

type fileslocalStarter struct {
	conf config.Config
}

var pathRepository string

// const pxPreviewDefault = 80
//var pxPreview = pxPreviewDefault
//var pxPreviewStr string

var _ starter.Operator = &fileslocalStarter{}

func (fc *fileslocalStarter) Name() string {
	return string(InterfaceKey)
}

func (fc *fileslocalStarter) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {

	fc.conf = conf

	confFileslocal, errs := conf.Fileslocal("", nil)
	if len(errs) < 1 {
		var ok bool
		if pathRepository, ok = confFileslocal["repository_path"]; !ok {
			errs = append(errs, fmt.Errorf("no repository_path in config.Fileslocal: %#v", confFileslocal))
		} else if _, err := os.Stat(pathRepository); err != nil {
			errs = append(errs, errors.Wrapf(err, "repository_path directory '%v' does not exist", pathRepository))
		}
	}

	//index, errs := config.ComponentIndex(indexPath, filelib.CurrentPath(), errs)
	//pxPreview_ := index.Options["pxPreview"]
	//if pxPreview_ != "" {
	//	var err error
	//	pxPreview, err = strconv.Atoi(pxPreview_)
	//	if err != nil {
	//		errs = append(errs, nil.Wrapf(err, "can't read 'pxPreview': %s", pxPreview_))
	//		pxPreview = pxPreviewDefault
	//	}
	//}
	//pxPreviewStr = strconv.Itoa(pxPreview)

	return nil, errs.Err()
}

func (fc *fileslocalStarter) Setup(conf config.Config, indexPath string, data map[string]string) error {
	err := os.MkdirAll(pathRepository, os.ModePerm)
	if err == nil {
		err = os.MkdirAll(pathRepository+"txt", os.ModePerm)
	}
	return err
}

func (fc *fileslocalStarter) Init() error {
	fileinfoOp, ok := joiner.Component(fileinfo.InterfaceKey).(fileinfo.Operator)
	if !ok {
		return errors.New("no group interface found :-(")
	}

	file, err := New(pathRepository, fileinfoOp)
	if err != nil {
		return errors.Wrap(err, "can't init files.Operator interface")
	}
	err = joiner.JoinInterface(file, files.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join files.Operator interface")
	}

	return nil
}
