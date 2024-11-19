package dynamictags

import (
	"os"
	"reflect"
)

const (
	ENV_TAG = "env"
)

type EnvTagConverter struct {
}

func NewEnvTagConverter() TagConverterer {
	return &EnvTagConverter{}
}

func (conv *EnvTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	val, isExists := os.LookupEnv(tag)
	return val, isExists, nil
}

func (conv EnvTagConverter) GetTag() string {
	return ENV_TAG
}
