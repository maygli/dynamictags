package dynamictags

import "reflect"

const (
	DEFAULT_TAG = "default"
)

type DefaultTagConverter struct {
}

// Set structure field with 'default' tag to value of the tag.
// Returns:
//   - Default tag converter.
func NewDefaultTagConverter() TagConverterer {
	return &DefaultTagConverter{}
}

// Returns conversion result.
// Parameters:
//   - tag tag value. This value already processed. All tokens like ${ENV_VARIABLE}
//         already replaced by dictionary value or environment variable value
//   - t structure field
//   - v value
//   - path json path to structure field
// Returns:
//   - Value which will set to structure field.
//   - Flag. If true value will be set. Otherwice it will be skiped
//   - error in case of error
func (conv *DefaultTagConverter) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error) {
	return tag, true, nil
}

// Returns converter tag.
// Returns:
//   - processed tag
func (conv DefaultTagConverter) GetTag() string {
	return DEFAULT_TAG
}
