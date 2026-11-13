package server_http

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/logger"
)

type ConfigStarter struct {
	Port        int    `yaml:"port"          json:"port"`
	NoHTTPS     bool   `yaml:"no_https"      json:"no_https"`
	KeyPath     string `yaml:"key_path"      json:"key_path"`
	TLSCertFile string `yaml:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile  string `yaml:"tls_key_file"  json:"tls_key_file"`
}

type Config struct {
	ConfigCommon
	EndpointsSettled
}

type ConfigCommon struct {
	Title   string
	Version string
	Host    string
	Port    string
	Prefix  string
}

// Swagger -----------------------------------------------------------------------------

type Swagger map[string]interface{}

func (c Config) SwaggerV2(isHTTPS bool) ([]byte, error) {
	paths := map[string]common.Map{} // map[string]map[string]map[string]interface{}{}

	for key, ep := range c.EndpointsSettled {

		path := c.Prefix + ep.Endpoint.PathTemplateBraced(ep.Path)
		method := strings.ToLower(ep.Endpoint.Method)

		epDescr := common.Map{
			"operationId": key,
			"tags":        ep.Tags,
		}

		if len(ep.Produces) >= 1 {
			epDescr["produces"] = ep.Produces
		} else {
			epDescr["produces"] = []string{"application/json"}
		}

		var parameters []interface{} // []map[string]interface{}

		for _, pp := range ep.Endpoint.PathParams {
			if len(pp) > 0 && pp[0] == '*' {
				pp = pp[1:]
			}

			parameters = append(
				parameters,
				common.Map{
					"in":          "path",
					"required":    true,
					"name":        pp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}
		for _, qp := range ep.Endpoint.QueryParams {
			parameters = append(
				parameters,
				common.Map{
					"in":          "query",
					"required":    false, // TODO!!!
					"name":        qp,
					"type":        "string",
					"description": "", // TODO!!!
				},
			)
		}

		if method == "post" {
			if len(ep.Endpoint.BodyParams) > 0 {
				parameters = append(parameters, ep.Endpoint.BodyParams)
			} else {
				parameters = append(parameters, common.Map{
					"in":       "body",
					"required": true,
					"name":     "body_item",
					"type":     "string",
				})
			}
		}

		if len(parameters) > 0 {
			epDescr["parameters"] = parameters
		}

		if epDescrPrev, ok := paths[path][method]; ok {
			return nil, fmt.Errorf("duplicate endpoint description (%s %s): \n%#v\nvs.\n%#v", method, path, epDescrPrev, epDescr)
		}
		if _, ok := paths[path]; ok { // pathPrev
			paths[path][method] = epDescr
		} else {
			paths[path] = common.Map{method: epDescr} // map[string]map[string]interface{}
		}
	}

	var schemes []string
	if isHTTPS {
		schemes = []string{"https", "http"}
	} else {
		schemes = []string{"http"}
	}

	swagger := Swagger{
		"swagger": "2.0",
		"info": map[string]string{
			"title":   c.Title,
			"version": c.Version,
		},
		// "basePath": c.Prefix,
		"schemes": schemes,
		"port":    c.Port,
		"paths":   paths,
	}

	return json.MarshalIndent(swagger, "", " ")
}

func (c Config) InitSwagger(isHTTPS bool, swaggerStaticFilePath string, l logger.Operator) error {
	//if c == nil {
	//	return nil
	//}

	swaggerJSON, err := c.SwaggerV2(isHTTPS)
	if err != nil {
		return err
	}
	if err = os.WriteFile(swaggerStaticFilePath, swaggerJSON, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0644): %s", swaggerStaticFilePath, swaggerJSON, err)
	}
	l.Infof("%d bytes are written into %s", len(swaggerJSON), swaggerStaticFilePath)

	return nil
}

//func (c *Config) SwaggerV2(isHTTPS bool) ([]byte, error) {
//
//  if c == nil {
//	  return nil, nil
//  }
//
//	paths := map[string]common.Map{} // map[string]map[string]map[string]interface{}{}
//
//	for key, ep := range c.EndpointsSettled {
//
//		path := c.Prefix + ep.PathTemplateBraced(ep.Path)
//		method := strings.ToLower(ep.Method)
//
//		epDescr := common.Map{
//			"operationId": key,
//			"tags":        ep.Tags,
//		}
//
//		if len(ep.Produces) >= 1 {
//			epDescr["produces"] = ep.Produces
//		} else {
//			epDescr["produces"] = []string{"application/json"}
//		}
//
//		var parameters []interface{} // []map[string]interface{}
//
//		for _, pp := range ep.PathParams {
//			if len(pp) > 0 && pp[0] == '*' {
//				pp = pp[1:]
//			}
//
//			parameters = append(
//				parameters,
//				common.Map{
//					"in":          "path",
//					"required":    true,
//					"name":        pp,
//					"type":        "string",
//					"description": "", // TODO!!!
//				},
//			)
//		}
//		for _, qp := range ep.QueryParams {
//			parameters = append(
//				parameters,
//				common.Map{
//					"in":          "query",
//					"required":    false, // TODO!!!
//					"name":        qp,
//					"type":        "string",
//					"description": "", // TODO!!!
//				},
//			)
//		}
//
//		if method == "post" {
//			if len(ep.BodyParams) > 0 {
//				parameters = append(parameters, ep.BodyParams)
//			} else {
//				parameters = append(parameters, common.Map{
//					"in":       "body",
//					"required": true,
//					"name":     "body_item",
//					"type":     "string",
//				})
//			}
//		}
//
//		if len(parameters) > 0 {
//			epDescr["parameters"] = parameters
//		}
//
//		if epDescrPrev, ok := paths[path][method]; ok {
//			return nil, fmt.Errorf("duplicate endpoint description (%s %s): \n%#v\nvs.\n%#v", method, path, epDescrPrev, epDescr)
//		}
//		if _, ok := paths[path]; ok { // pathPrev
//			paths[path][method] = epDescr
//		} else {
//			paths[path] = common.Map{method: epDescr} // map[string]map[string]interface{}
//		}
//	}
//
//	var schemes []string
//	if isHTTPS {
//		schemes = []string{"https", "http"}
//	} else {
//		schemes = []string{"http"}
//	}
//
//	swagger := Swagger{
//		"swagger": "2.0",
//		"info": map[string]string{
//			"title":   c.Title,
//			"version": c.Version,
//		},
//		// "basePath": c.Prefix,
//		"schemes": schemes,
//		"port":    c.Port,
//		"paths":   paths,
//	}
//
//	return json.MarshalIndent(swagger, "", " ")
//}
//
