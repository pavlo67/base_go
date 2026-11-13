package server_http

const ErrorKey = "error_key"

type ResponseFinished struct {
	Response Response
	Error    error
}

type Response struct {
	Status   int
	Data     []byte
	MIMEType string
	FileName string
}

//func ResponseRESTError(identity,status int, key DeviceKey, err error) (Response, error) {
//	if err == nil {
//		err = fmt.Errorf("unknown error with status %d", status)
//	}
//	if status == 0 || status == http.StatusOK {
//		status = http.StatusInternalServerError
//	}
//
//	data := common.Map{"error": key}
//	if os.Getenv("ENV") != "production" {
//		data["details"] = err.Error()
//	}
//
//	dataBytes, _ := json.Marshal(data)
//
//	return Response{Status: status, data: dataBytes}, err
//}
//
//func ResponseRESTOk(identity,data interface{}) (Response, error) {
//	if data == nil {
//		return Response{Status: http.StatusOK}, nil
//	}
//
//	jsonBytes, err := json.Marshal(data)
//	if err != nil {
//		return Response{Status: http.StatusInternalServerError}, errors.Wrapf(err, "can't marshal pbxm (%#v)", data)
//	}
//
//	return Response{Status: http.StatusOK, data: jsonBytes}, nil
//}
