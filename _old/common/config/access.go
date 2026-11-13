package config

import (
	"strconv"
	"strings"
)

// DEPRECATED!!! use url.URL instead
type Access struct {
	Proto   string
	Host    string
	Port    int
	User    string
	Pass    string
	Path    string
	Options string
}

func (access *Access) URL() string {
	if access == nil {
		return ""
	}

	url := strings.TrimSpace(access.Host)
	if access.Port > 0 {
		url += ":" + strconv.Itoa(access.Port)
	}

	path := strings.TrimSpace(access.Path)
	if len(path) > 0 && path[0] != '/' {
		url += "/" + path
	} else {
		url += path
	}

	if url != "" {
		proto := strings.TrimSpace(access.Proto)

		// TODO!!! remove the kostyl
		if proto == "" {
			proto = "http://"
		}

		url = proto + url
	}

	return url
}
