package server_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

const (
	CORSAllowHeaders     = "authorization,content-type"
	CORSAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	CORSAllowOrigin      = "*"
	CORSAllowCredentials = "true"
)

// REST -------------------------------------------------------------------------------------

type RESTDataMessage struct {
	Info     string `json:"info,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

// Redirect ----------------------------------------------------------------------------------

func Redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func ResponseRESTError(status int, err error, req *http.Request) (Response, error) {
	commonErr := errors.CommonError(err)

	key := commonErr.Key()
	data := common.Map{ErrorKey: key}

	if status == 0 || status == http.StatusOK {
		if key == common.NoCredsKey || key == common.InvalidCredsKey {
			status = http.StatusUnauthorized
		} else if key == common.NoUserKey || key == common.NoRightsKey {
			status = http.StatusForbidden
		} else if status == 0 || status == http.StatusOK {
			status = http.StatusInternalServerError

		} else {
			status = http.StatusInternalServerError
		}
	}

	//if !strlib.In(s.secretENVsToLower, strings.ToLower(os.Getenv("ENV"))) {
	//	data["details"] = commonErr.Error()
	//}

	jsonBytes, errJSON := json.Marshal(data)
	if errJSON != nil {
		commonErr = commonErr.Append(fmt.Errorf("marshalling error data (%#v): %s", data, errJSON))
	}

	if req != nil {
		commonErr = commonErr.Append(fmt.Errorf("on %s %s", req.Method, req.URL))
	}

	return Response{Status: status, Data: jsonBytes}, commonErr
}

func ResponseRESTOk(status int, data interface{}, req *http.Request) (Response, error) {
	if status == 0 {
		status = http.StatusOK
	}

	if data == nil {
		return Response{Status: status}, nil
	}

	var dataBytes []byte

	switch v := data.(type) {
	case []byte:
		dataBytes = v
	case *[]byte:
		if v != nil {
			dataBytes = *v
		}
	case string:
		dataBytes = []byte(v)
	case *string:
		if v != nil {
			dataBytes = []byte(*v)
		}
	default:
		var err error
		if dataBytes, err = json.Marshal(data); err != nil {
			if req != nil {
				err = fmt.Errorf("on %s %s: can't marshal json (%#v), got %s", req.Method, req.URL, data, err)
			} else {
				err = fmt.Errorf("can't marshal json (%#v): %s", data, err)
			}

			return Response{Status: http.StatusInternalServerError}, err
		}
	}

	return Response{Status: status, Data: dataBytes}, nil
}
