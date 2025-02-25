package dynamictags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	INT_PREFIX_KEY             = "INT_PREFIX"
	INT_PREFIX_VAL             = "2"
	UINT_PREFIX_KEY            = "UINT_PREFIX"
	UINT_PREFIX_VAL            = "12"
	FLOAT_PREFIX_KEY           = "FLOAT_PREFIX"
	FLOAT_PREFIX_VAL           = "67."
	STRING_PREFIX_KEY          = "STRING_PREFIX"
	STRING_PREFIX_VAL          = "String"
	BOOL_PREFIX_KEY            = "BOOL_PREFIX"
	BOOL_PREFIX_VAL            = "True"
	INTERNAL_INT_PREFIX_KEY    = "INTERNAL_INT_PREFIX"
	INTERNAL_INT_PREFIX_VAL    = "72"
	INTERNAL_FLOAT_PREFIX_KEY  = "INTERNAL_FLOAT_PREFIX"
	INTERNAL_FLOAT_PREFIX_VAL  = "52."
	INTERNAL_STRING_PREFIX_KEY = "INTERNAL_STRING_PREFIX"
	INTERNAL_STRING_PREFIX_VAL = "Str"
)

const (
	EXPECTED_INT             = 21
	EXPECTED_UINT            = 12123
	EXPECTED_FLOAT           = 67.78
	EXPECTED_STRING          = "StringTest"
	EXPECTED_NO_TAG          = ""
	EXPECTED_BOOL            = true
	EXPECTED_BLACK_NO_BLACK  = "BlackList"
	EXPECTED_BLACK_BLACK     = ""
	EXPECTED_INTERNAL_INT    = 72456789
	EXPECTED_INTERNAL_FLOAT  = 52.99
	EXPECTED_INTERNAL_STRING = "StrInternalString"
	EXPECTED_SLICE_SIZE      = 3
	EXPECTED_SLICE_VAL0      = "one"
	EXPECTED_SLICE_VAL1      = "two"
	EXPECTED_SLICE_VAL2      = "free"
)

type TestInternalStruct struct {
	IntData    int64   `default:"${INTERNAL_INT_PREFIX}456789"`
	FloatData  float64 `default:"${INTERNAL_FLOAT_PREFIX}99"`
	StringData string  `default:"${INTERNAL_STRING_PREFIX}InternalString"`
}

type TestStruct struct {
	IntData       int8     `default:"${INT_PREFIX}1"`
	UIntData      uint16   `default:"${UINT_PREFIX}123"`
	FloatData     float32  `default:"${FLOAT_PREFIX}78"`
	StringData    string   `default:"${STRING_PREFIX}Test"`
	BoolData      bool     `default:"${BOOL_PREFIX}"`
	BlackListData string   `default:"BlackList"`
	Slice         []string `default:"one,two,free"`
	unField       string   `default:"unField"`
	NoTagData     string
	IntStruct     TestInternalStruct
}

func verifyResult(t *testing.T, res TestStruct, isWithBalckList bool) {
	assert.Equal(t, int8(EXPECTED_INT), res.IntData)
	assert.Equal(t, uint16(EXPECTED_UINT), res.UIntData)
	assert.Equal(t, float32(EXPECTED_FLOAT), res.FloatData)
	assert.Equal(t, EXPECTED_STRING, res.StringData)
	assert.Equal(t, EXPECTED_NO_TAG, res.NoTagData)
	assert.Equal(t, EXPECTED_BOOL, res.BoolData)
	assert.Equal(t, EXPECTED_SLICE_SIZE, len(res.Slice))
	assert.Equal(t, EXPECTED_SLICE_VAL0, res.Slice[0])
	assert.Equal(t, EXPECTED_SLICE_VAL1, res.Slice[1])
	assert.Equal(t, EXPECTED_SLICE_VAL2, res.Slice[2])
	assert.Equal(t, "", res.unField)
	if isWithBalckList {
		assert.Equal(t, res.BlackListData, EXPECTED_BLACK_BLACK)
	} else {
		assert.Equal(t, res.BlackListData, EXPECTED_BLACK_NO_BLACK)
	}
	assert.Equal(t, int64(EXPECTED_INTERNAL_INT), res.IntStruct.IntData)
	assert.Equal(t, EXPECTED_INTERNAL_FLOAT, res.IntStruct.FloatData)
	assert.Equal(t, EXPECTED_INTERNAL_STRING, res.IntStruct.StringData)
}

func TestCreateDeafultProcessor(t *testing.T) {
	processor := NewDefaultProcessor()
	assert.NotNil(t, processor)
	assert.Equal(t, 1, len(processor.converters))
	conv, ok := processor.converters[0].(*DefaultTagConverter)
	assert.True(t, ok)
	assert.NotNil(t, conv)
}

func TestDefaultProcessorSimple(t *testing.T) {
	defaultProcessor := NewDefaultProcessor()
	defaultProcessor.SetDictionaryValue(INT_PREFIX_KEY, INT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(UINT_PREFIX_KEY, UINT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(FLOAT_PREFIX_KEY, FLOAT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(STRING_PREFIX_KEY, STRING_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(BOOL_PREFIX_KEY, BOOL_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_INT_PREFIX_KEY, INTERNAL_INT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_FLOAT_PREFIX_KEY, INTERNAL_FLOAT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_STRING_PREFIX_KEY, INTERNAL_STRING_PREFIX_VAL)
	testStruct := TestStruct{}
	err := defaultProcessor.Process(&testStruct, nil)
	assert.NoError(t, err)
	verifyResult(t, testStruct, false)
}

func TestDefaultProcessorBlackList(t *testing.T) {
	defaultProcessor := NewDefaultProcessor()
	defaultProcessor.SetDictionaryValue(INT_PREFIX_KEY, INT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(UINT_PREFIX_KEY, UINT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(FLOAT_PREFIX_KEY, FLOAT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(STRING_PREFIX_KEY, STRING_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(BOOL_PREFIX_KEY, BOOL_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_INT_PREFIX_KEY, INTERNAL_INT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_FLOAT_PREFIX_KEY, INTERNAL_FLOAT_PREFIX_VAL)
	defaultProcessor.SetDictionaryValue(INTERNAL_STRING_PREFIX_KEY, INTERNAL_STRING_PREFIX_VAL)
	blackList := make([]string, 0)
	blackList = append(blackList, "$.BlackListData")
	testStruct := TestStruct{}
	err := defaultProcessor.Process(&testStruct, blackList)
	assert.NoError(t, err)
	verifyResult(t, testStruct, true)
}
