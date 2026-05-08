package filelib

import (
	"path/filepath"
	"strings"
)

func Join(basePath, path string) string {
	if path = strings.TrimSpace(path); path == "" {
		return basePath
	} else if filepath.VolumeName(path) != "" {
		return path
	} else if filepath.ToSlash(path)[0] == '/' {
		return path
	}

	return filepath.Join(basePath, path)

}
