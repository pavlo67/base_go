package httplib

import (
	"testing"

	"github.com/pavlo67/common/common/logger/logger_test"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	//_, cfg, l := apps.PrepareTests(t, "../../_envs/", "test", "request.log")
	//require.NotNil(t, cfg)

	url := "http://google.com/"
	method := "GET"
	var responseData []byte
	l := logger_test.New(t, "", "", false, nil)

	err := Request(nil, url, method, nil, nil, &responseData, l, "")
	require.NoError(t, err)

	// t.Logf("%s")
}
