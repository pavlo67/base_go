package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"
)

func SetCreds(creds auth.Creds) http.Header {
	var header http.Header

	if jwt := creds[auth.CredsJWT]; jwt != "" {
		if header == nil {
			header = http.Header{}
		}
		header.Add("Authorization", jwt)
	} else if token := creds[auth.CredsToken]; token != "" {
		if header == nil {
			header = http.Header{}
		}
		header.Add("Authorization", token)
	}

	return header
}
