package loader_http

import (
	"net/url"
	"strings"
)

type Priority func(string) int

func PriorityDefault(uriDefault string, noDoubts bool) Priority {
	var domainDefault string

	u, _ := url.Parse(uriDefault)
	if u != nil {
		domainDefault = strings.ToLower(u.Hostname())
	}

	if domainDefault == "" && !noDoubts {
		return func(_ string) int {
			return -1
		}
	}

	return func(uri string) int {
		var domain string

		u, _ := url.Parse(uri)
		if u != nil {
			domain = strings.ToLower(u.Hostname())
		}

		if domain == domainDefault {
			return 1
		}

		return -1
	}
}
