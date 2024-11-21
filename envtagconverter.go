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

// Set structure field with 'env' tag to value of environment variable.
// Returns:
//   - Environment variable converter.
func NewEnvTagConverter() TagConverterer {
	return &EnvTagConverter{}
}

// Returns conversion result.
// Parameters:
//   - tag tag value. This value already processed. All tokens like ${ENV_VARIABLE}
//     already replaced by dictionary value or environment variable value
//   - t structure field
//   - v value
//   - path json path to structure field
//
// Returns:
//   - Value which will set to structure field.
//   - Flag. If true value will be set. Otherwice it will be skiped
//   - error in case of error
func (conv *EnvTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	val, isExists := os.LookupEnv(tag)
	return val, isExists, nil
}

// Returns converter tag.
// Returns:
//   - processed tag
func (conv EnvTagConverter) GetTag() string {
	return ENV_TAG
}
