package serialization

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var MarshalerYAML = MarshalerStruct{yaml.Marshal, yaml.Unmarshal}
var MarshalerJSON = MarshalerStruct{json.Marshal, json.Unmarshal}

// MarshalerStruct ----------------------------------------------------------------------------------------

type Marshal func(v interface{}) ([]byte, error)
type Unmarshal func(data []byte, v interface{}) error

var _ Marshaler = &MarshalerStruct{}

type MarshalerStruct struct {
	marshal   Marshal
	unmarshal Unmarshal
}

func (cs MarshalerStruct) Marshal(v interface{}) ([]byte, error) {
	return cs.marshal(v)
}

func (cs MarshalerStruct) Unmarshal(data []byte, v interface{}) error {
	return cs.unmarshal(data, v)
}
