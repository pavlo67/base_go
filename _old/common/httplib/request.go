package httplib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server_http"
)

const bodyLogLimit = 2048

const onRequest = "on httplib.Request()"

type ResponseBinary struct {
	MIMEType string
	Data     []byte
}

var reBinary = regexp.MustCompile(`^image/`)

func Request(client *http.Client, serverURL, method string, header http.Header, requestData, responseData interface{}, l logger.Operator, logFile string) error {
	if client == nil {
		client = &http.Client{}
	}

	var err error
	//for _, doReAuth := range reAuthTries {

	// start of single try

	var requestBody []byte
	var requestBodyReader io.Reader

	if requestData != nil {
		switch v := requestData.(type) {
		case []byte:
			requestBody = v
		case *[]byte:
			requestBody = *v
		case string:
			requestBody = []byte(v)
		case *string:
			requestBody = []byte(*v)
		default:
			if requestBody, err = json.Marshal(requestData); err != nil {
				return fmt.Errorf(onRequest+": can't marshal request responseData (%#v): %s", requestData, err)
			}
		}

		// must be checked for nil instead direct write
		// the external for GET requests expected nil body, but nil-requestData after json.Marshal return not empty responseData

		requestBodyReader = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequest(method, serverURL, requestBodyReader)
	if err != nil || req == nil {
		logger.LogRequest(method, serverURL, nil, requestBody, nil, nil, err, 0, logFile)
		return fmt.Errorf(onRequest+": can't create request %s %s, got %#v, %s", method, serverURL, req, err)
	} else if req.Body != nil {
		defer Close(req.Body, client, nil)
	}
	req.Header = header
	var responseBody []byte

	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer Close(resp.Body, client, nil)
	}

	if err != nil {
		var statusCode int
		var responseHeaders http.Header
		if resp != nil {
			statusCode = resp.StatusCode
			responseHeaders = resp.Header
			responseBody, _ = ioutil.ReadAll(resp.Body)
		}

		logger.LogRequest(method, serverURL, req.Header, requestBody, responseHeaders, responseBody, err, statusCode, logFile)
		return fmt.Errorf(onRequest+": can't %s %s, got %s", method, serverURL, err)
	}

	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Error("can't read response body: ", err)
	}

	var responseBodyToLog []byte

	if !reBinary.MatchString(resp.Header.Get("Content-Type")) {
		responseBodyToLog = responseBody
	}
	logger.LogRequest(method, serverURL, req.Header, requestBody, resp.Header, responseBodyToLog, err, resp.StatusCode, logFile)
	if err != nil {
		return fmt.Errorf(onRequest+": can't read body from %s %s, got %s", method, serverURL, err)
	}

	if resp.StatusCode == http.StatusUnauthorized { // && doReAuth
		//if identity.Token = reAuthJWT(*identity); identity.Token != "" {
		//	continue
		//}
	}

	if resp.StatusCode != http.StatusOK {
		// TODO!!! be careful writing server_http handlers, http.StatusOK is the only success code accepted here

		if len(responseBody) > bodyLogLimit {
			responseBody = responseBody[:bodyLogLimit]
		}

		var data common.Map
		if err = json.Unmarshal(responseBody, &data); err != nil {
			if len(responseBody) > bodyLogLimit {
				responseBody = responseBody[:bodyLogLimit]
			}
			return fmt.Errorf(onRequest+": can't unmarshal body from %s %s: status = %d, body = %s, got %s", method, serverURL, resp.StatusCode, responseBody, err)
		}

		errCommon := fmt.Sprintf("can't %s %s: status = %d, body = %s", method, serverURL, resp.StatusCode, responseBody)
		if data["error"] != nil {
			data["error"] = errors.CommonError(data["error"], errCommon)
		} else {
			data["error"] = errCommon
		}
		errorKey := errors.Key(data.StringDefault(server_http.ErrorKey, ""))
		return errors.CommonError(errorKey, data)
	}

	if responseData != nil {
		switch v := responseData.(type) {
		case *[]byte:
			if v != nil {
				*v = responseBody
			}
		case *string:
			if v != nil {
				*v = string(responseBody)
			}
		case *ResponseBinary:
			v.MIMEType = resp.Header.Get("Content-Type")
			v.Data = responseBody
		default:
			if err = json.Unmarshal(responseBody, responseData); err != nil {
				if len(responseBody) > bodyLogLimit {
					responseBody = responseBody[:bodyLogLimit]
				}
				return fmt.Errorf(onRequest+": can't unmarshal body from %s %s %s, got %s", method, serverURL, responseBody, err)
			}
		}
	}

	//	break // end of each try means the end of all tries if something other wasn't managed before
	//}

	return nil
}

//func RequestJSON(method, url string, data []byte, headers map[string]string) (common.Map, error) {
//	client := &http.Client{}
//
//	req, err := http.NewRequest(method, url, bytes.NewReader(data))
//	if err != nil {
//		return nil, err
//	}
//
//	for k, v := range headers {
//		req.Header.Add(k, v)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	// log.Printf("%s", body)
//
//	result := common.Map{}
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return result, errors.Wrapf(err, "can't unmarsal: %s", body)
//	}
//
//	return result, nil
//}
