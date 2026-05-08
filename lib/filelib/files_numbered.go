package filelib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type NumberedFile struct {
	I    int
	Path string
}

func NumberedFilesSequence(dir, filenameRegexp string, allowLeaks bool) ([]NumberedFile, error) {
	reFilename, err := regexp.Compile(filenameRegexp)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fns []NumberedFile
	nextNum := -1

	for _, file := range files {
		if matches := reFilename.FindStringSubmatch(file.Name()); len(matches) == 2 {
			num, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("converting num (%v) from filename (%s) got: %s", matches, file.Name(), err)
			} else if nextNum == -1 {
				nextNum = num
			} else if num != nextNum && !allowLeaks {
				return nil, fmt.Errorf("frames leak: after %d got %d", nextNum-1, num)
			}
			fns = append(fns, NumberedFile{num, filepath.Join(dir, file.Name())})
			nextNum++
		}
	}

	return fns, nil
}
