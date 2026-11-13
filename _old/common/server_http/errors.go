package server_http

import (
	"fmt"
	"net/http"
)

func On(req *http.Request) string {
	if req == nil {
		return ""
	}
	return fmt.Sprintf("ERROR on %s %s: ", req.Method, req.URL)
}
