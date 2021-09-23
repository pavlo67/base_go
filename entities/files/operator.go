package files

import (
	"os"

	"github.com/pavlo67/data/types"
)

type Operator interface {
	Save(path, newFilePattern string, data []byte) (string, error)
	Read(path string) ([]byte, error)
	Remove(path string) error
	List(path string, depth int) (Items, error)
	Stat(path string, depth int) (*types.Stat01, error)
}

type Items []types.File01

func (fis Items) Append(basePath string, info os.FileInfo) (Items, error) {
	path := info.Name()

	//if len(path) <= len(basePath) {
	//	return nil, fmt.Errorf("wrong path (%s) on basePath = '%s'", path, basePath)
	//}

	if info.IsDir() {
		if path != "" && path[len(path)-1] != '/' {
			path += "/"
		}
		fis = append(fis, types.File01{
			Path: path,
			// Path:      path[len(basePath):],
			IsDir:     true,
			CreatedAt: info.ModTime(),
		})
	} else {
		fis = append(fis, types.File01{
			Path: path,
			// Path:      path[len(basePath):],
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return fis, nil
}
