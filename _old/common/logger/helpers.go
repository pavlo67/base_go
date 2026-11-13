package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pavlo67/common/common/filelib"
)

// const onLogTimestamp = "on logger.LogTimestamp()"

//func LogTimestamp(l Operator) {
//	l.Info("\n\ncreated at: " + time.Now().Format(time.RFC3339))
//
//	//f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	//if err != nil {
//	//	Log(l, err.Error(), true, onLogTimestamp)
//	//}
//	//defer f.Close()
//	//if _, err := f.WriteString(message); err != nil {
//	//	Log(l, err.Error(), true, onLogTimestamp)
//	//}
//}

var MaxLoggedDataLength = 20000

func LogRequest(method, path string, reqHeaders http.Header, reqBody []byte, respHeaders http.Header, respBody []byte, bodyErr error,
	status int, logFile string) {
	if len(reqBody) > MaxLoggedDataLength {
		reqBody = reqBody[:MaxLoggedDataLength]
	}
	if len(respBody) > MaxLoggedDataLength {
		respBody = respBody[:MaxLoggedDataLength]
	}

	ts := time.Now().Format(time.RFC3339)

	data := fmt.Sprintf("\n\n%s %s %s\nheaders: %#v\nbody: %s\nresponse: %d %s %#v", method, path, ts,
		reqHeaders, reqBody, status, respBody, respHeaders)
	if bodyErr != nil {
		data = fmt.Sprintf("\nERROR: %s", bodyErr) + data
	}

	fmt.Println(data)
	if logFile != "" {
		filelib.AppendFile(logFile, []byte(data+"\n"))
	}
}

//func Log(l Operator, data string, isError bool, message string) {
//	data = message + ": " + data
//	if l == nil {
//		if isError {
//			data = "ERROR: " + data
//		}
//		log.Print(data)
//	} else if isError {
//		l.Error(data)
//	} else {
//		l.Info(data)
//	}
//}
//
//func LogIntoFile(logfile string, l Operator, data, message string) {
//
//	data += "\n"
//	Log(l, data, false, message)
//
//	if logfile != "" {
//		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//		if err != nil {
//			Log(l, err.Error(), true, message)
//		}
//		defer f.Close()
//		if _, err := f.WriteString(data); err != nil {
//			Log(l, err.Error(), true, message)
//		}
//	}
//}
//
//func FatalIntoFile(logfile string, l Operator, data, message string) {
//
//	data += "\n"
//	Log(l, data, true, message)
//
//	if logfile != "" {
//		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//		if err != nil {
//			Log(l, err.Error(), true, message)
//		}
//		defer f.Close()
//		if _, err := f.WriteString(data); err != nil {
//			Log(l, err.Error(), true, message)
//		}
//	}
//
//	os.Exit(1)
//}
