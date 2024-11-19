package dynamictags

import "reflect"

// Interface for tag converter
// Tag Converter get tag string and returns value which will be set to processed
// structure
type TagConverterer interface {
	// This function is called for all simple (like int8,int16, etc.) structure fields.
	// Should returns value corresponded to tag.
	// Parameters:
	//   - tag processed tag. To this function passed processed tag. I.e. all substring like
	//     '${KEY}' already replaced by dictionary or environment variable value
	//   - t reflect structure field
	//   - v reflect structure value
	//   - path json path to structure filed (like '$.InternalStructure.Data1')
	//
	// Returns:
	//   - result string (for example value of environment variable with name specified in tag)
	//   - if 'false' result value will not be set (for example if the environment variable is not exists)
	//   - error in case of error. Processing will be interrupted and the error value will returned
	GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (any, bool, error)

	// Return processed tag
	// Returns:
	// - tag
	GetTag() string
}
