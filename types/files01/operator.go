package files01

import (
	"os"
	"time"
)

type Operator interface {
	Save(path, newFilePattern string, data []byte) (string, error)
	Read(path string) ([]byte, error)
	Remove(path string) error
	List(path string, depth int) (Items, error)
	Stat(path string, depth int) (*Item, error)
}

type Item struct {
	Path string
	// Name      string
	IsDir     bool
	Size      int64
	CreatedAt time.Time
}

type Items []Item

func (fis Items) Append(basePath string, info os.FileInfo) (Items, error) {
	path := info.Name()

	//if len(path) <= len(basePath) {
	//	return nil, fmt.Errorf("wrong path (%s) on basePath = '%s'", path, basePath)
	//}

	if info.IsDir() {
		if path != "" && path[len(path)-1] != '/' {
			path += "/"
		}
		fis = append(fis, Item{
			Path: path,
			// Path:      path[len(basePath):],
			IsDir:     true,
			CreatedAt: info.ModTime(),
		})
	} else {
		fis = append(fis, Item{
			Path: path,
			// Path:      path[len(basePath):],
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return fis, nil
}
