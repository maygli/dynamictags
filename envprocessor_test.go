package dynamictags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnv(t *testing.T) {
	processor := NewEnvProcessor()
	assert.NotNil(t, processor)
	assert.Equal(t, 1, len(processor.converters))
	conv, ok := processor.converters[0].(*EnvTagConverter)
	assert.True(t, ok)
	assert.NotNil(t, conv)
}
