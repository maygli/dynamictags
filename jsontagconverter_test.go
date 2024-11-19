package dynamictags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonTagConverter(t *testing.T) {
	converter := NewJsonTagConverter()
	assert.NotNil(t, converter)
}
