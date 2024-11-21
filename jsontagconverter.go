package dynamictags

import (
	"reflect"

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
	resMap, ok := conv.jsonData.(map[string]interface{})
	if !ok {
		return "", false, nil
	}
	val, ok := resMap[tag]
	return val, ok, nil
}

func (conv JsonTagConverter) GetTag() string {
	return JSON_TAG
}
