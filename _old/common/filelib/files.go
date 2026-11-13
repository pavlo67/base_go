package filelib

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func CurrentFile(removeExt bool) string {
	if _, filename, _, ok := runtime.Caller(1); ok {
		if removeExt {
			filenamePOSIX := reBackslash.ReplaceAllString(filename, "/")
			partes := strings.Split(filenamePOSIX, "/")
			if len(partes) < 1 {
				return ""
			}
			partes[len(partes)-1] = reExt.ReplaceAllString(partes[len(partes)-1], "")
			return strings.Join(partes, "/")
		}
		return filename
	}
	return ""
}

func CorrectFileName(name string) string {
	name = rePoint.ReplaceAllLiteralString(name, "_")
	name = reSpecials.ReplaceAllLiteralString(name, "_")
	return name
}

const onCopyFiles = "on filelib.CopyFiles()"

func CopyFiles(path, targetPath string, reStr string, removeOrigin bool) error {
	var re *regexp.Regexp
	var err error

	if len(reStr) > 0 {
		re, err = regexp.Compile(reStr)
		if err != nil {
			return fmt.Errorf("wrong reStr: '%s'"+onCopyFiles, reStr)
		}
	}

	targetPath, err = Dir(targetPath)
	if err != nil {
		return fmt.Errorf("targetPath is wrong: %s / "+onCopyFiles, err)
	}

	dirEntries, err := List(path, re, false, true)
	if err != nil {
		return fmt.Errorf("%s / "+onCopyFiles, err)
	}

	for _, dirEntry := range dirEntries {
		if removeOrigin {
			err = os.Rename(dirEntry, filepath.Join(targetPath, filepath.Base(dirEntry)))
		} else {
			err = CopyFile(dirEntry, filepath.Join(targetPath, filepath.Base(dirEntry)))
		}

		if err != nil {
			return fmt.Errorf("%s / "+onCopyFiles, err)
		}

	}

	return nil
}

func BackupFile(fileName string) error {
	dir := path.Dir(fileName) + "/"

	backupName := dir + path.Base(fileName) + "." + time.Now().Format(time.RFC3339)[:10] + ".bak"

	_, err := os.Stat(backupName)

	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = CopyFile(fileName, backupName)
		return errors.Wrapf(err, "on copying '%s' to '%s'", fileName, backupName)
	}

	return errors.Wrapf(err, "on copying '%s' to '%s'", fileName, backupName)
}

type Convert func(data []byte) ([]byte, error)

func CopyFileConverted(src, dst string, perm fs.FileMode, convert Convert) error {
	// TODO??? check if src and dst are the same

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if convert != nil {
		data, err = convert(data)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(dst, data, perm)
}

func CopyFile(src, dst string) error {
	// TODO??? check if src and dst are the same

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err = in.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing %s: %s", src, err)
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err = out.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing %s: %s", dst, err)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func RemoveContents(dir string) error {
	//d, err := os.Open(dir)
	//if err != nil {
	//	return err
	//}
	//defer d.Close()
	//names, err := d.Readdirnames(-1)
	//if err != nil {
	//	return err
	//}

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err = os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

const onRemoveFile = "on filelib.RemoveFile()"

func RemoveFile(path string) error {
	if fileExists, _ := FileExistsAny(path); fileExists {
		if err := os.Remove(path); err != nil {
			return errors.Wrap(err, onRemoveFile)
		}
	}
	return nil
}
