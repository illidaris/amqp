package amqp

import (
	"encoding/json"
	"errors"
)

type Encoder interface {
	GetContentType() string
	GetEncoding() string
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type JSONEncoder struct{}

func (e JSONEncoder) GetContentType() string {
	return "application/json"
}
func (e JSONEncoder) GetEncoding() string {
	return "utf8"
}
func (e JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (e JSONEncoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type defaultEncoder struct{}

func (e defaultEncoder) GetContentType() string {
	return "text/plain"
}
func (e defaultEncoder) GetEncoding() string {
	return "utf8"
}
func (e defaultEncoder) Encode(v interface{}) ([]byte, error) {
	if str, ok := v.(string); ok {
		return []byte(str), nil
	}
	return nil, errors.New("value is not string")
}

func (e defaultEncoder) Decode(data []byte, v interface{}) error {
	value := string(data)
	v = &value
	return nil
}
