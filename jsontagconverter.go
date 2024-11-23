package dynamictags

import (
	"reflect"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

const (
	JSON_TAG = "json"
)

type JsonTagConverter struct {
	jsonData any
}

func NewJsonTagConverter(content any, rootPath string) (TagConverterer, error) {
	conv := JsonTagConverter{}
	var err error
	conv.jsonData, err = jsonpath.Get(rootPath, content)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (conv *JsonTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	jsonPath := path
	if !strings.HasPrefix(tag, "$") {
		// If path is relative
		jsonPath = path + "." + tag
	}
	data, err := jsonpath.Get(jsonPath, conv.jsonData)
	return data, err == nil, nil
}

func (conv JsonTagConverter) GetTag() string {
	return JSON_TAG
}
