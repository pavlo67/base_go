package node_crud_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &demoStarter{}
}

var _ starter.Operator = &demoStarter{}

type demoStarter struct{}

// --------------------------------------------------------------------------

var l logger.Operator

func (ds *demoStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ds *demoStarter) Prepare(cfg *config.Config, options common.Map) error {
	return nil
}

// Swagger-UI sorts interface sections due to the first their path occurrences, so:
// 1. unauthorized   /auth/...
// 2. admin          /front/...

// TODO!!! keep in mind that EndpointsConfig key and corresponding .HandlerKey not necessarily are the same, they can be defined different

var ServerConfig = server_http.Config{
	ConfigCommon: server_http.ConfigCommon{
		Title:   "NODE_CRUD REST API",
		Version: "0.0.1",
	},
	EndpointsSettled: server_http.EndpointsSettled{
		auth.IntefaceKeyAuthenticate: {Path: "/auth", Tags: []string{"unauthorized"}},
		auth.IntefaceKeySetCreds:     {Path: "/set_creds", Tags: []string{"unauthorized"}},
	},
}

const PrefixREST = "/backend"

func (ds *demoStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if srvOp == nil {
		return fmt.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvPort, isHTTPS := srvOp.Addr()
	restStaticPath := filelib.CurrentPath() + "./api-docs/"

	if err := ServerConfig.CompleteWithJoiner(joinerOp, "", srvPort, PrefixREST); err != nil {
		return err
	}
	if err := ServerConfig.HandleEndpoints(srvOp, l); err != nil {
		return err
	}

	if err := ServerConfig.InitSwagger(isHTTPS, restStaticPath+"swagger.json", l); err != nil {
		return err
	}
	if err := srvOp.HandleFiles("rest_static", PrefixREST+"/api-docs/*filepath", server_http.StaticPath{LocalPath: restStaticPath}); err != nil {
		return err
	}

	WG.Add(1)

	// TODO!!! customize it
	// if isHTTPS {
	//	go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
	// }

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	return nil
}
