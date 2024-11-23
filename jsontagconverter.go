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
	data, err := jsonpath.Get(path+"."+tag, conv.jsonData)
	return data, err == nil, nil
}

func (conv JsonTagConverter) GetTag() string {
	return JSON_TAG
}
