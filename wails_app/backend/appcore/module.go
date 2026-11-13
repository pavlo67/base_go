package appcore

import (
	"encoding/json"
	"net/http"
)

type Module interface {
	Name() string
	RegisterHTTP(mux *http.ServeMux)
}

type Registry struct {
	modules []Module
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) Register(module Module) {
	r.modules = append(r.modules, module)
}

func (r *Registry) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("/api/modules", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		names := make([]string, 0, len(r.modules))
		for _, module := range r.modules {
			names = append(names, module.Name())
		}
		WriteJSON(w, http.StatusOK, map[string]any{"modules": names})
	})

	for _, module := range r.modules {
		module.RegisterHTTP(mux)
	}
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WithLocalDevCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
