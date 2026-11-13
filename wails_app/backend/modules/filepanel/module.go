package filepanel

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"wails-vue-go/backend/appcore"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/entities/files/files_fs"
	"github.com/pavlo67/base_go/lib/logger"
)

const Name = "filepanel"

type Module struct {
	startPath string
	fsRoot    string
	op        files.Operator
}

type Entry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Size     uint64 `json:"size"`
	MimeType string `json:"mimeType"`
	MTime    string `json:"mtime,omitempty"`
	IsParent bool   `json:"isParent"`
}

type ListResponse struct {
	Module    string  `json:"module"`
	Path      string  `json:"path"`
	Parent    string  `json:"parent,omitempty"`
	CanGoUp   bool    `json:"canGoUp"`
	Entries   []Entry `json:"entries"`
	StartPath string  `json:"startPath"`
	FSRoot    string  `json:"fsRoot"`
}

func NewFromStartup(l logger.Operator) (*Module, error) {
	startPath, err := startupPath()
	if err != nil {
		return nil, err
	}
	fsRoot := filesystemRoot(startPath)

	op, _, err := files_fs.New(fsRoot, l)
	if err != nil {
		return nil, err
	}

	return &Module{startPath: startPath, fsRoot: fsRoot, op: op}, nil
}

func (m *Module) Name() string { return Name }

func (m *Module) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("/api/filepanel/config", m.handleConfig)
	mux.HandleFunc("/api/filepanel/list", m.handleList)
}

func (m *Module) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appcore.WriteJSON(w, http.StatusOK, map[string]any{
		"module":    Name,
		"startPath": m.startPath,
		"fsRoot":    m.fsRoot,
	})
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimSpace(r.URL.Query().Get("path"))
	if path == "" {
		path = m.startPath
	}

	path, err := filepath.Abs(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isInside(m.fsRoot, path) {
		http.Error(w, "path is outside filesystem root", http.StatusBadRequest)
		return
	}

	items, err := m.op.List(path, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entries := directChildren(path, items)
	parent := filepath.Dir(path)
	canGoUp := samePath(path, parent) == false && isInside(m.fsRoot, parent)
	if canGoUp {
		entries = append([]Entry{{
			Name:     "..",
			Path:     filepath.ToSlash(parent),
			IsDir:    true,
			MimeType: files_fs.MimeTypeDirectory,
			IsParent: true,
		}}, entries...)
	} else {
		parent = ""
	}

	appcore.WriteJSON(w, http.StatusOK, ListResponse{
		Module:    Name,
		Path:      filepath.ToSlash(path),
		Parent:    filepath.ToSlash(parent),
		CanGoUp:   canGoUp,
		Entries:   entries,
		StartPath: m.startPath,
		FSRoot:    m.fsRoot,
	})
}

func directChildren(dir string, items []files.Item) []Entry {
	entries := make([]Entry, 0, len(items))
	for _, item := range items {
		itemPath := filepath.Clean(filepath.FromSlash(item.Path))
		if filepath.Dir(itemPath) != filepath.Clean(dir) {
			continue
		}

		entries = append(entries, Entry{
			Name:     filepath.Base(itemPath),
			Path:     filepath.ToSlash(itemPath),
			IsDir:    item.IsDir,
			Size:     item.Size,
			MimeType: item.MimeType,
			MTime:    item.MTime.Format("2006-01-02 15:04:05"),
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries
}

func startupPath() (string, error) {
	if fromEnv := strings.TrimSpace(os.Getenv("WAILS_APP_START_DIR")); fromEnv != "" {
		return filepath.Abs(fromEnv)
	}

	if len(os.Args) > 1 {
		if arg := strings.TrimSpace(os.Args[1]); arg != "" {
			return filepath.Abs(arg)
		}
	}

	return os.Getwd()
}

func filesystemRoot(path string) string {
	clean := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		volume := filepath.VolumeName(clean)
		if volume != "" {
			return volume + string(os.PathSeparator)
		}
		return `C:\`
	}
	return string(os.PathSeparator)
}

func isInside(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	if samePath(root, path) {
		return true
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func samePath(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

func formatSize(size uint64) string {
	return strconv.FormatUint(size, 10)
}

func (e Entry) String() string {
	if e.IsDir {
		return fmt.Sprintf("[%s]", e.Path)
	}
	return fmt.Sprintf("%s (%s bytes)", e.Path, formatSize(e.Size))
}
