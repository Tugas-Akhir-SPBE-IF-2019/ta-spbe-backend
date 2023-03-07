package jsonmanipulator

import "encoding/json"

type Client interface {
	Marshal(v any) ([]byte, error)
}

type simpleJSONManipulator struct{}

func NewSimpleJSONManipulator() (Client, error) {
	return &simpleJSONManipulator{}, nil
}

func (jsonEC simpleJSONManipulator) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
