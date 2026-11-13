package server_http

import (
	"encoding/json"

	"github.com/pavlo67/common/common/joiner"
)

type EndpointDescription struct {
	InternalKey joiner.InterfaceKey `json:",omitempty"`
	Method      string              `json:",omitempty"`
	PathParams  []string            `json:",omitempty"`
	QueryParams []string            `json:",omitempty"`
	BodyParams  json.RawMessage     `json:",omitempty"`
}

type EndpointKey = joiner.InterfaceKey
type EndpointsSettled map[EndpointKey]EndpointSettled

type EndpointSettled struct {
	Path     string
	Tags     []string `json:",omitempty"`
	Produces []string `json:",omitempty"`
	*Endpoint
}

type Endpoint struct {
	EndpointDescription
	WorkerHTTP
}
