package files01_fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pavlo67/data/elements/ns"
	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/data/entities/files01"
)

var _ files01.Operator = &filesFS{}

type filesFS struct {
	basePath string
}

const onNew = "on filesFS.New(): "

func New(basePath string) (files01.Operator, db.Cleaner, error) {
	filesOp := filesFS{}

	var err error
	if filesOp.basePath, err = filelib.Dir(basePath); err != nil || filesOp.basePath == "" {
		return nil, nil, fmt.Errorf(onNew+": creating base path '%s' got %s", basePath, err)
	}

	// l.Infof("%s --> %s", basePath, filesOp.basePath)

	return &filesOp, &filesOp, nil
}

const onSave = "on filesFS.Save()"

func (filesOp *filesFS) Save(path, newFilePattern string, data []byte) (string, error) {
	basePath := filesOp.basePath
	path = basePath + path

	var err error
	var dirPath string
	var file *os.File

	// TODO!!! check if dirPath doesn't contain "/../"
	if newFilePattern != "" {
		if dirPath, err = filelib.Dir(path); err != nil {
			return "", errors.Wrapf(err, onSave+": wrong path (%s)", path)
		}
		if file, err = ioutil.TempFile(dirPath, newFilePattern); err != nil {
			return "", errors.Wrapf(err, onSave+": can't ioutil.TempFile(%s, %s)", dirPath, newFilePattern)
		}
	} else {
		var filename string
		dirPath, filename = filepath.Dir(path), filepath.Base(path)

		dirPath, err = filelib.Dir(dirPath)

		if file, err = os.OpenFile(dirPath+filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			return "", errors.Wrapf(err, onSave+": can't os.OpenFile(%s, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)", dirPath+filename)
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			l.Errorf(onSave+": on file.Close() got %s", err)
		}
	}()

	filename := strings.ReplaceAll(file.Name(), "/./", "/")

	if len(filename) <= len(basePath) {
		return "", fmt.Errorf(onSave+": wrong filename (%s) on basePath = '%s'", filename, basePath)
	}

	if _, err = file.Write(data); err != nil {
		return "", errors.Wrapf(err, onSave+": can't file.Write(%s)", file.Name())
	}

	return filename[len(basePath):], nil
}

const onRead = "on filesFS.Read()"

func (filesOp *filesFS) Read(path string) ([]byte, error) {
	filePath := filesOp.basePath + path

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't ioutil.ReadFile(%s)", filePath)
	}

	return data, nil
}

const onRemove = "on filesFS.Remove()"

func (filesOp *filesFS) Remove(path string) error {
	filePath := filesOp.basePath + path

	if err := os.Remove(filePath); err != nil {
		return errors.Wrapf(err, onRemove+": can't os.Remove(%s)", filePath)
	}

	return nil
}

const onList = "on filesFS.Items()"

func (filesOp *filesFS) List(path string, depth int) (files01.Items, error) {
	filePath := filesOp.basePath + path

	var filesInfo files01.Items

	if depth == 0 {
		fis, err := ioutil.ReadDir(filePath)

		if err != nil {
			return nil, errors.Wrapf(err, onList+": can't ioutil.ReadDir(%s)", filePath)
		}

		for _, fi := range fis {
			filesInfo, err = filesInfo.Append("", fi) // basePath
			if err != nil {
				return nil, errors.Wrap(err, onList)
			}
		}

		return filesInfo, nil
	}

	// TODO: process depth > 0 more thoroughly here
	err := filepath.Walk(filePath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filesInfo, err = filesInfo.Append("", fi) // basePath
		if err != nil {
			return errors.Wrap(err, onList)
		}

		return nil
	})

	return filesInfo, err
}

const onStat = "on filesFS.Stat()"

func (filesOp *filesFS) Stat(path string, depth int) (*entities.Stat, error) {
	filePath := filesOp.basePath + path

	fi, err := os.Stat(filePath)
	if err != nil {
		//if os.IsNotExist(err) {
		//	return nil, nil
		//}
		return nil, errors.Wrapf(err, onStat+": can't  os.Stat(%s)", filePath)
	}

	filesInfo, err := files01.Items{}.Append("", fi) // basePath
	if err != nil || len(filesInfo) != 1 {
		return nil, fmt.Errorf(onStat+": got %#v / %s", filesInfo, err)
	}

	fileInfo := filesInfo[0]
	stat := entities.Stat{
		NSS:       ns.NSS(fileInfo.Path),
		TotalSize: fileInfo.Size,
		CreatedAt: fileInfo.CreatedAt,
	}

	if depth != 0 && fileInfo.IsDir {
		// TODO: process depth > 0 more thoroughly here
		err = filepath.Walk(filePath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			stat.ChldCount++
			stat.TotalSize += fi.Size()
			return nil
		})
	}

	return &stat, err

}
