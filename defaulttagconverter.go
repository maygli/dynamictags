package dynamictags

import "reflect"

const (
	DEFAULT_TAG = "default"
)

type DefaultTagConverter struct {
}

func NewDefaultTagConverter() TagConverterer {
	return &DefaultTagConverter{}
}

func (conv *DefaultTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	return tag, true, nil
}

func (conv DefaultTagConverter) GetTag() string {
	return DEFAULT_TAG
}
