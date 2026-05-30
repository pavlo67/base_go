package files_fs

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/lib/db"
	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/filelib"
	"github.com/pavlo67/base_go/lib/logger"
)

const MimeTypeDirectory = "inode/directory"
const MimeTypeSymlink = "inode/symlink"

var _ files.Operator = &filesFS{}
var _ db.Operator = &filesFS{}

type filesFS struct {
	root string
}

var l logger.Operator

const onNew = "on files_fs.New():"

func New(root string, l_ logger.Operator) (files.Operator, db.Operator, error) {
	if l_ == nil {
		return nil, nil, errors.New("", onNew+" l_ == nil")
	}
	l = l_

	if root == "" {
		return nil, nil, errors.New("", onNew+" root is empty")
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}
	if err = os.MkdirAll(rootAbs, os.ModePerm); err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	return &filesFS{root: rootAbs}, &filesFS{root: rootAbs}, nil
}

const onSave = "on files_fs.Save():"

func (op *filesFS) Save(data files.Data) error {
	path, err := op.fullPath(data.Path)
	if err != nil {
		return errors.Wrap(err, onSave)
	}

	if data.IsDir {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.Wrap(err, onSave)
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return errors.Wrap(err, onSave)
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return errors.Wrap(err, onSave)
		}
		if err = f.Truncate(int64(data.Size)); err != nil {
			if closeErr := f.Close(); closeErr != nil {
				l.Errorf(onSave+": on file.Close(): %v", closeErr)
			}
			return errors.Wrap(err, onSave)
		}
		if err = f.Close(); err != nil {
			return errors.Wrap(err, onSave)
		}
	}

	if !data.MTime.IsZero() {
		if err = os.Chtimes(path, time.Now(), data.MTime); err != nil {
			return errors.Wrap(err, onSave)
		}
	}

	return nil
}

const onRead = "on files_fs.Read():"

func (op *filesFS) Read(path string) (*files.Item, error) {
	fullPath, err := op.fullPath(path)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	info, err := os.Lstat(fullPath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	itemPath, err := op.dataPath(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return op.item(itemPath, info), nil
}

const onRemove = "on files_fs.Remove():"

func (op *filesFS) Remove(path string, forceRecursion bool) error {
	fullPath, err := op.fullPath(path)
	if err != nil {
		return errors.Wrap(err, onRemove)
	}
	if fullPath == op.root {
		return fmt.Errorf(onRemove + " removing root is not allowed")
	}

	if forceRecursion {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}
	if os.IsNotExist(err) {
		return nil
	}

	return errors.Wrap(err, onRemove)
}

const onList = "on files_fs.List():"

func (op *filesFS) List(path string, depth int) ([]files.Item, error) {
	if depth < 0 {
		return nil, errors.New("", onList+" depth must be >= 0")
	}

	fullPath, err := op.fullPath(path)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	rootInfo, err := os.Lstat(fullPath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, onList)
	}
	if !rootInfo.IsDir() {
		return nil, fmt.Errorf(onList+" path is not a directory: %s", path)
	}

	var items []files.Item
	err = filepath.WalkDir(fullPath, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if current == fullPath {
			return nil
		}

		relDepth, err := relDepth(fullPath, current)
		if err != nil {
			return err
		}
		if depth > 0 && relDepth > depth {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		itemPath, err := op.dataPath(current)
		if err != nil {
			return err
		}
		items = append(items, *op.item(itemPath, info))

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	return items, nil
}

func (op *filesFS) fullPath(path string) (string, error) {
	return filelib.FullPath(op.root, path)
}

// db.Operator --------------------------------------------------------------------------------

const onCreate = "on files_fs.Create():"

func (op *filesFS) Create(_ *sql.DB) error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf(onCreate+" allowed only when ENV=TEST, but env = %s", env)
	}

	return errors.Wrap(os.MkdirAll(op.root, os.ModePerm), onCreate)
}

const onClean = "on files_fs.Clean():"

func (op *filesFS) Clean(opts interface{}) error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf(onClean+" allowed only when ENV=TEST, but env = %s", env)
	}

	path, ok := opts.(string)
	if !ok {
		return fmt.Errorf(onClean+" expects directory path as string, got %T", opts)
	}
	if path == "" {
		return fmt.Errorf(onClean + " expects non-empty directory path")
	}

	fullPath, err := op.fullPath(path)
	if err != nil {
		return errors.Wrap(err, onClean)
	}

	if fullPath == op.root {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return errors.Wrap(err, onClean)
		}
		for _, entry := range entries {
			if err = os.RemoveAll(filepath.Join(fullPath, entry.Name())); err != nil {
				return errors.Wrap(err, onClean)
			}
		}
		return nil
	}

	if err = os.RemoveAll(fullPath); err != nil {
		return errors.Wrap(err, onClean)
	}
	return errors.Wrap(os.MkdirAll(fullPath, os.ModePerm), onClean)
}
