package dynamictags

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	EXPECTED_TAG = "default"
	TEST_STRING  = "test"
)

func TestDefaultTagConverter(t *testing.T) {
	conv := NewDefaultTagConverter()
	assert.NotNil(t, conv)
	assert.Equal(t, EXPECTED_TAG, conv.GetTag())
	val, isSet, err := conv.GetSimpleValue(TEST_STRING, reflect.StructField{}, reflect.Value{}, "")
	assert.Equal(t, TEST_STRING, val)
	assert.True(t, isSet)
	assert.NoError(t, err)
}
