package httplib

import (
	"io"
	"log"
	"net/http"

	"github.com/pavlo67/common/common/logger"
)

func Close(readCloser io.Closer, client *http.Client, l logger.Operator) {
	if err := readCloser.Close(); err != nil {
		if l != nil {
			l.Error(err)
		} else {
			log.Print(err)
		}
	}

	if client != nil {
		client.CloseIdleConnections()
	}
}
