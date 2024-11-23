package dynamictags

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	JSON_DATA = `{
		"root" : {
			"testcfg1" : {
				"val1" : 123,
				"val2" : true,
				"val3" : "testdata"
			}
		}
}
`
	TEST_PATH         = "$.root.testcfg1"
	INCORRECT_PATH    = "$.data"
	EXPECTED_JSON_TAG = "json"
)

func TestJsonTagConverter(t *testing.T) {
	var res any
	err := json.Unmarshal([]byte(JSON_DATA), &res)
	assert.NoError(t, err)

	// Case 1 correct create tag converter
	conv, err := NewJsonTagConverter(res, TEST_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	assert.Equal(t, EXPECTED_JSON_TAG, conv.GetTag())
	// Get existed value
	val, ok, err := conv.GetSimpleValue("val1", reflect.StructField{}, reflect.Value{}, "$")
	assert.NotNil(t, val)
	assert.NoError(t, err)
	assert.True(t, ok)
	// Get non existed value
	val, ok, err = conv.GetSimpleValue("val1567", reflect.StructField{}, reflect.Value{}, "$")
	assert.Nil(t, val)
	assert.NoError(t, err)
	assert.False(t, ok)

	// Case 2 json data is nil
	conv, err = NewJsonTagConverter(nil, TEST_PATH)
	assert.Error(t, err)
	assert.Nil(t, conv)

	// Case 3 incorrect path
	conv, err = NewJsonTagConverter(nil, INCORRECT_PATH)
	assert.Error(t, err)
	assert.Nil(t, conv)
}
