package files_fs

import (
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pavlo67/base_go/entities/files"
)

func (op *filesFS) item(path string, info os.FileInfo) *files.Item {
	mtime := info.ModTime().UTC()

	return &files.Item{
		Data: files.Data{
			Path:     filepath.ToSlash(path),
			IsDir:    info.IsDir(),
			Size:     uint64(info.Size()),
			MTime:    mtime,
			CRC:      nil,
			MimeType: mimeType(path, info),
		},
		CreatedAt: time.Time{},
		UpdatedAt: mtime,
	}
}

func mimeType(path string, info os.FileInfo) string {
	if info.Mode()&os.ModeSymlink != 0 {
		return MimeTypeSymlink
	}
	if info.IsDir() {
		return MimeTypeDirectory
	}
	if mimeType := mime.TypeByExtension(filepath.Ext(path)); mimeType != "" {
		return mimeType
	}
	return "application/octet-stream"
}

func (op *filesFS) dataPath(path string) (string, error) {
	rel, err := filepath.Rel(op.root, path)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return filepath.ToSlash(op.root), nil
	}
	return filepath.ToSlash(path), nil
}

func relDepth(basePath, path string) (int, error) {
	rel, err := filepath.Rel(basePath, path)
	if err != nil {
		return 0, err
	}
	if rel == "." {
		return 0, nil
	}
	return len(strings.Split(filepath.ToSlash(rel), "/")), nil
}
