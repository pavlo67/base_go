package files_http_server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/logger"
)

var l logger.Operator

const onRegister = "on files_http_server.Register():"

func Register(mux *http.ServeMux, prefix string, filesOp files.Operator, l_ logger.Operator) error {
	if mux == nil {
		return errors.New("", onRegister+" mux == nil")
	}
	if filesOp == nil {
		return errors.New("", onRegister+" filesOp == nil")
	}
	if l_ == nil {
		return errors.New("", onRegister+" l_ == nil")
	}
	l = l_

	prefix = strings.TrimRight(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		return errors.New("", onRegister+" prefix is empty")
	}

	mux.HandleFunc(prefix+"/save", handleSave(filesOp))
	mux.HandleFunc(prefix+"/read", handleRead(filesOp))
	mux.HandleFunc(prefix+"/remove", handleRemove(filesOp))
	mux.HandleFunc(prefix+"/list", handleList(filesOp))

	return nil
}

// handlers ------------------------------------------------------------------------------------

const onSave = "on files_http_server.Save():"

func handleSave(op files.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}

		var data files.Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			writeError(w, http.StatusBadRequest, errors.Wrap(err, onSave))
			return
		}
		if err := op.Save(data); err != nil {
			writeError(w, http.StatusInternalServerError, errors.Wrap(err, onSave))
			return
		}
		writeJSON(w, http.StatusOK, okResponse{OK: true})
	}
}

const onRead = "on files_http_server.Read():"

func handleRead(op files.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}

		path := r.URL.Query().Get("path")
		item, err := op.Read(path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, errors.Wrap(err, onRead))
			return
		}
		writeJSON(w, http.StatusOK, readResponse{Item: item})
	}
}

const onRemove = "on files_http_server.Remove():"

func handleRemove(op files.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}

		path := r.URL.Query().Get("path")
		force, err := strconv.ParseBool(r.URL.Query().Get("force"))
		if err != nil {
			writeError(w, http.StatusBadRequest, errors.Wrap(err, onRemove))
			return
		}
		if err = op.Remove(path, force); err != nil {
			writeError(w, http.StatusInternalServerError, errors.Wrap(err, onRemove))
			return
		}
		writeJSON(w, http.StatusOK, okResponse{OK: true})
	}
}

const onList = "on files_http_server.List():"

func handleList(op files.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}

		depth, err := strconv.Atoi(r.URL.Query().Get("depth"))
		if err != nil {
			writeError(w, http.StatusBadRequest, errors.Wrap(err, onList))
			return
		}
		items, err := op.List(r.URL.Query().Get("path"), depth)
		if err != nil {
			writeError(w, http.StatusInternalServerError, errors.Wrap(err, onList))
			return
		}
		writeJSON(w, http.StatusOK, listResponse{Items: items})
	}
}

// response helpers ----------------------------------------------------------------------------

type okResponse struct {
	OK bool `json:"ok"`
}

type readResponse struct {
	Item *files.Item `json:"item"`
}

type listResponse struct {
	Items []files.Item `json:"items"`
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeError(w, http.StatusMethodNotAllowed, errors.New("", "method not allowed"))
	return false
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		l.Errorf("on files_http_server.writeJSON(): %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{Error: err.Error()})
}
