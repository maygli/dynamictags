package dynamictags

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/PaesslerAG/jsonpath"
)

const (
	JSON_TAG = "json"
)

type JsonTagConverter struct {
}

func NewJsonTagConverter() TagConverterer {
	val := interface{}(nil)
	err := json.Unmarshal([]byte("{\"message\":{\"test\":122}}"), &val)
	if err != nil {
		return nil
	}
	jsonmsg, err := jsonpath.Get("$.message.test", val)
	if err != nil {
		return nil
	}
	jsonmap, ok := jsonmsg.(map[string]interface{})
	if !ok {
		return nil
	}
	fmt.Println(jsonmap)
	return &EnvTagConverter{}
}

func (conv *JsonTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	return "", false, nil
}

func (conv JsonTagConverter) GetTag() string {
	return JSON_TAG
}
