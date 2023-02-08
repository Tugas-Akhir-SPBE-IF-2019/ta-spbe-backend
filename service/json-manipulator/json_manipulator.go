package jsonmanipulator

import "encoding/json"

type EncoderDecoder struct{}

func (jsonEC EncoderDecoder) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
