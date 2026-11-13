package server_http

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Config -----------------------------------------------------------------------------------

func (c *Config) EP(endpointKey EndpointKey, params []string, createFullURL bool) (string, string, error) {
	if c == nil {
		return "", "", nil
	}

	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return "", "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if len(ep.PathParams) != len(params) {
		return "", "", fmt.Errorf("wrong params list (%#v) for endpoint (%s / %#v)", params, endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	for i, param := range params {
		if param == "" {
			return "", "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(param)
	}

	return ep.Method, urlStr, nil
}

var rePathParam = regexp.MustCompile(":[^/]+")

func (ed EndpointDescription) PathTemplate(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ed.PathParams) < 1 {

		// TODO!!! be careful with serverPath like ".../*something"
		// if serverPath[len(serverPath)-1] != '/' {
		//	serverPath += "/"
		// }
		return serverPath

	} else if serverPath[len(serverPath)-1] == '/' {
		serverPath = serverPath[:len(serverPath)-1]
	}
	var pathParams []string
	for _, pp := range ed.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp)
		} else {
			pathParams = append(pathParams, ":"+pp)
		}

	}

	return serverPath + "/" + strings.Join(pathParams, "/")
}

func (ed EndpointDescription) PathTemplateBraced(serverPath string) string {
	if len(serverPath) == 0 || serverPath[0] != '/' {
		serverPath = "/" + serverPath
	}

	if len(ed.PathParams) < 1 {
		return serverPath
	}

	var pathParams []string
	for _, pp := range ed.PathParams {
		if len(pp) > 0 && pp[0] == '*' {
			pathParams = append(pathParams, pp[1:])
		} else {
			pathParams = append(pathParams, pp)
		}

	}

	return serverPath + "/{" + strings.Join(ed.PathParams, "}/{") + "}"
}

//func (ep Endpoint) PathWithParams(params ...string) string {
//	matches := rePathParam.FindAllStringSubmatchIndex(ep.Path, -1)
//
//	numMatches := len(matches)
//	if len(params) < numMatches {
//		numMatches = len(params)
//	}
//
//	path := ep.Path
//	for nm := numMatches - 1; nm >= 0; nm-- {
//		path = path[:matches[nm][0]] + url.PathEscape(strings.ReplaceTags(params[nm], "/", "%2F", -1)) + path[matches[nm][1]:]
//	}
//
//	return path
//}

// this trick allows to prevent run-time errors with wrong endpoint parameters number
// using CheckGet...() functions we move parameter number checks to initiation stage

type Get1 func(string) (string, error)
type Get2 func(string, string) (string, error)
type Get3 func(string, string, string) (string, error)
type Get4 func(string, string, string, string) (string, error)

func CheckGet0(c Config, endpointKey EndpointKey, createFullURL bool) (string, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return "", fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	if strings.ToUpper(ep.Method) != "GET" {
		return "", fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	return urlStr, nil
}

func CheckGet1(c Config, endpointKey EndpointKey, createFullURL bool) (Get1, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1 string) (string, error) {
		p1 = strings.TrimSpace(p1)
		if p1 == "" {
			return "", fmt.Errorf("empty param %s for endpoint (%s / %#v)", p1, endpointKey, ep)
		}
		urlStr += "/" + url.PathEscape(p1)
		return urlStr, nil
	}, nil
}

func CheckGet2(c Config, endpointKey EndpointKey, createFullURL bool) (Get2, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2 string) (string, error) {
		params := [2]string{p1, p2}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet3(c Config, endpointKey EndpointKey, createFullURL bool) (Get3, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3 string) (string, error) {
		params := [3]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}

func CheckGet4(c Config, endpointKey EndpointKey, createFullURL bool) (Get4, error) {
	ep, ok := c.EndpointsSettled[endpointKey]
	if !ok {
		return nil, fmt.Errorf("no endpoint with key '%s'", endpointKey)
	}

	var urlStr string
	if createFullURL {
		urlStr = c.Host
		if c.Port = strings.TrimSpace(c.Port); c.Port != "" {
			urlStr += ":" + c.Port
		}
	}
	urlStr += c.Prefix + ep.Path

	if strings.ToUpper(ep.Method) != "GET" {
		return nil, fmt.Errorf("wrong endpoint.Method with key '%s': %#v", endpointKey, ep)
	}

	return func(p1, p2, p3, p4 string) (string, error) {
		params := [4]string{p1, p2, p3}
		for i, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				return "", fmt.Errorf("empty param %d in list (%#v) for endpoint (%s / %#v)", i, params, endpointKey, ep)
			}
			urlStr += "/" + url.PathEscape(param)
		}
		return urlStr, nil
	}, nil
}
