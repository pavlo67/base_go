package server_http_jschmhr

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server_http"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	port        int
	tlsCertFile string
	tlsKeyFile  string

	onRequest server_http.OnRequestMiddleware

	secretENVsToLower []string
}

func New(port int, tlsCertFile, tlsKeyFile string, secretENVs []string) (server_http.Operator, error) {
	if port <= 0 {
		return nil, fmt.Errorf("on server_http_jschmhr.New(): wrong port = %d", port)
	}

	var secretENVsToLower []string
	for _, secretENV := range secretENVs {
		secretENVsToLower = append(secretENVsToLower, strings.ToLower(secretENV))
	}

	router := httprouter.New()

	return &serverHTTPJschmhr{
		httpServer: &http.Server{
			Handler:        router,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,

			// ReadHeaderTimeout:            60 * time.Second,
			// IdleTimeout:                  60 * time.Second,
			// DisableGeneralOptionsHandler: false,
		},
		httpServeMux: router,
		port:         port,
		tlsCertFile:  tlsCertFile,
		tlsKeyFile:   tlsKeyFile,

		secretENVsToLower: secretENVsToLower,
	}, nil
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *serverHTTPJschmhr) Start() error {
	if s == nil {
		return errors.New("no serverOp to start")
	}

	s.httpServer.Addr = ":" + strconv.Itoa(s.port)
	l.Info("Server is starting on address ", s.httpServer.Addr)

	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		return s.httpServer.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
	}

	return s.httpServer.ListenAndServe()
}

func (s *serverHTTPJschmhr) Shutdown(ctx context.Context) error {
	if s == nil {
		return errors.New("no serverOp to shutting down")
	}

	return s.Shutdown(ctx)
}

func (s *serverHTTPJschmhr) Addr() (port int, https bool) {
	return s.port, s.tlsCertFile != "" && s.tlsKeyFile != ""
}

//func (s *serverHTTPJschmhr) ServerHTTP() *http.Server {
//	return s.httpServer
//}

const onHandleMiddleware = "on serverHTTPJschmhr.HandleMiddleware()"

func (s *serverHTTPJschmhr) HandleMiddleware(onRequest server_http.OnRequestMiddleware) error {
	if s.onRequest != nil && onRequest != nil {
		return fmt.Errorf(onHandleMiddleware + ": can't add middlware twice")
	}

	s.onRequest = onRequest
	return nil
}

const onHandleEndpoint = "on serverHTTPJschmhr.HandleEndpoint()"

func (s *serverHTTPJschmhr) HandleEndpoint(key server_http.EndpointKey, serverPath string, endpoint server_http.Endpoint) error {

	method := strings.ToUpper(endpoint.Method)
	path := endpoint.PathTemplate(serverPath)

	if endpoint.WorkerHTTP == nil {
		return errors.New(onHandleEndpoint + ": " + method + ": " + path + "\t!!! NULL workerHTTP ISN'T DISPATCHED !!!")
	}

	s.HandleOptions(key, path)

	handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		var identity *auth.Identity
		if s.onRequest != nil {
			var err error
			if identity, err = s.onRequest.Identity(r); err != nil {
				l.Error(err)
			}
		}

		var params server_http.PathParams
		if len(paramsHR) > 0 {
			params = server_http.PathParams{}
			for _, p := range paramsHR {
				params[p.Key] = p.Value
			}
		}

		responseData, err := endpoint.WorkerHTTP(s, r, params, identity)
		if err != nil {
			l.Warn(err)
		}

		if responseData.MIMEType != "" {
			w.Header().Set("Content-Type", responseData.MIMEType)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Content-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err = w.Write(responseData.Data); err != nil {
			l.Errorf("can't write response (%s): %s", serverPath, err)
		}
	}

	l.Infof("%-10s: %s %s", key, method, path)
	switch method {
	case "GET":
		s.httpServeMux.GET(path, handler)
	case "POST":
		s.httpServeMux.POST(path, handler)
	case "PUT":
		s.httpServeMux.PUT(path, handler)
	case "DELETE":
		s.httpServeMux.DELETE(path, handler)
	default:
		return fmt.Errorf(onHandleEndpoint+": method (%s) isn't supported", method)
	}

	return nil
}

func (s *serverHTTPJschmhr) HandleOptions(key server_http.EndpointKey, serverPath string) {
	//if strlib.In(s.handledOptions, serverPath) {
	//	//l.Infof("- %#v", s.handledOptions)
	//	return
	//}

	s.httpServeMux.OPTIONS(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l.Infof("%-10s: OPTIONS %s", key, serverPath)
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
	})

	//s.handledOptions = append(s.handledOptions, serverPath)
}

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *serverHTTPJschmhr) HandleFiles(key server_http.EndpointKey, serverPath string, staticPath server_http.StaticPath) error {
	l.Infof("%-10s: FILES %s <-- %s", key, serverPath, staticPath.LocalPath)

	// TODO: check localPath

	if staticPath.MIMEType == nil {
		// TODO!!! CORS

		s.httpServeMux.ServeFiles(serverPath, http.Dir(staticPath.LocalPath))
		return nil
	}

	s.HandleOptions(key, serverPath)

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		if staticPath.MIMEType != nil && *staticPath.MIMEType != "" {
			w.Header().Set("Content-Type", *staticPath.MIMEType)
		}

		OpenFile, err := os.Open(staticPath.LocalPath + "/" + p.ByName("filepath"))
		defer OpenFile.Close()
		if err != nil {
			l.Error(err)
		} else {
			io.Copy(w, OpenFile)
		}

		//if mimeType != nil {
		//}
		//fileServer.ServeHTTP(w, r)
	})

	return nil
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.ExportID.PathWithParams, nil)
// if err != nil {
//	l.ErrStr("can't read MIMEType for file: ", localPath+"/"+r.ExportID.PathWithParams, err)
// }

//func (s *serverHTTPJschmhr) HandleGetString(serverRoute, str string, mimeType *string) {
//	s.handleFunc("GET", serverRoute, func(w http.ResponseWriter, r *http.Request, params httprouter.Content) {
//		if mimeType != nil {
//			// "application/javascript"
//			w.Header().Set("Content-Type", *mimeType)
//		}
//		w.Write([]byte(str))
//	})
//}
