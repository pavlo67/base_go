package filelib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/pavlo67/common/common/errors"
)

const onSearch = "on filelib.Search()"

func Search(path string, re regexp.Regexp, getFirst bool) ([]string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, onSearch)
	}

	names := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		names[i] = dirEntry.Name()
	}

	slices.Sort(names)
	// sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

	var matches []string
	for _, name := range names {
		if matches = re.FindStringSubmatch(name); getFirst && len(matches) > 0 {
			return matches, nil
		}
	}

	return matches, nil
}

const onList = "on filelib.List()"

func List(path string, re *regexp.Regexp, getDirs, getFiles bool) ([]string, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	var names []string
	for _, dirEntry := range dirEntries {
		if (dirEntry.IsDir() && !getDirs) || (!dirEntry.IsDir() && !getFiles) {
			continue
		}
		names = append(names, dirEntry.Name())
	}

	slices.Sort(names)

	var namesListed []string
	for _, name := range names {
		if re == nil || re.MatchString(name) {
			namesListed = append(namesListed, filepath.Join(path, name))
		}
	}

	return namesListed, nil
}

const onFileExists = "on filelib.FileExists()"

func FileExists(path string, isDir bool) (bool, error) {
	fileInfo, _ := os.Stat(path)
	if fileInfo == nil {
		return false, nil
	}

	if fileInfo.IsDir() {
		if isDir {
			return true, nil
		}
		return false, fmt.Errorf("%s is not a directory / "+onFileExists, path)
	}

	if isDir {
		return false, fmt.Errorf("%s is a directory / "+onFileExists, path)
	}

	return true, nil
}

const onFileExistsAny = "on filelib.FileExistsAny()"

func FileExistsAny(path string) (exists, isDir bool) {
	fileInfo, _ := os.Stat(path)
	if fileInfo == nil {
		return false, false
	}

	return true, fileInfo.IsDir()
}
