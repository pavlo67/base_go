package filelib

import (
	"golang.org/x/sys/unix"

	"github.com/pavlo67/common/common/errors"
)

const onFreeSpace = "on filelib.FreeSpace()"

func FreeSpace(path string) (uint64, error) {
	// TODO!!! check if UNIX

	//path, err := os.Getwd()

	var stat unix.Statfs_t
	if err := unix.Statfs(path, &stat); err != nil {
		return 0, errors.Wrap(err, onFreeSpace)
	}

	// Available blocks * size per block = available space in bytes
	return stat.Bavail * uint64(stat.Bsize), nil
}
