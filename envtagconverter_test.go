package dynamictags

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	EXPECTED_ENV_TAG = "env"
	TEST_ENV_STRING  = "TEST_ENV"
	TEST_ENV_VALUE   = "VALUE"
)

func TestEnvConverter(t *testing.T) {
	conv := NewEnvTagConverter()
	assert.NotNil(t, conv)
	assert.Equal(t, EXPECTED_ENV_TAG, conv.GetTag())
	// Case 1 no environmant variable defined
	val, isSet, err := conv.GetSimpleValue(TEST_ENV_STRING, reflect.StructField{}, reflect.Value{}, "")
	assert.Empty(t, val)
	assert.False(t, isSet)
	assert.NoError(t, err)
	// Case 2 environment variable is defined
	os.Setenv(TEST_ENV_STRING, TEST_ENV_VALUE)
	val, isSet, err = conv.GetSimpleValue(TEST_ENV_STRING, reflect.StructField{}, reflect.Value{}, "")
	assert.Equal(t, TEST_ENV_VALUE, val)
	assert.True(t, isSet)
	assert.NoError(t, err)
}
