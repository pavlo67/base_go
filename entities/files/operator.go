package files

import (
	"os"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/elements/crud"
)

type Operator interface {
	Save(path, newFilePattern string, data []byte) (string, error)
	Read(path string) ([]byte, error)
	Remove(path string) error
	List(path string, depth int) (Items, error)
	Stat(path string, depth int) (*crud.Stat, error)
}

type Items []entities.File01

func (fis Items) Append(basePath string, info os.FileInfo) (Items, error) {
	path := info.Name()

	//if len(path) <= len(basePath) {
	//	return nil, fmt.Errorf("wrong path (%s) on basePath = '%s'", path, basePath)
	//}

	modTime := info.ModTime()

	if info.IsDir() {
		if path != "" && path[len(path)-1] != '/' {
			path += "/"
		}
		fis = append(fis, entities.File01{
			Path: path,
			// Path:      path[len(basePath):],
			IsDir:     true,
			UpdatedAt: &modTime,
		})
	} else {
		fis = append(fis, entities.File01{
			Path: path,
			// Path:      path[len(basePath):],
			Size:      info.Size(),
			UpdatedAt: &modTime,
		})
	}

	return fis, nil
}
