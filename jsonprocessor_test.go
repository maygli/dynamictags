package dynamictags

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	JSON_INT_PREFIX_KEY    = "JSON_INT_PREFIX"
	JSON_UINT_PREFIX_KEY   = "JSON_UINT_PREFIX"
	JSON_FLOAT_PREFIX_KEY  = "JSON_FLOAT_PREFIX"
	JSON_STRING_PREFIX_KEY = "JSON_STRING_PREFIX"
	JSON_BOOL_PREFIX_KEY   = "JSON_BOOL_PREFIX"
	JSON_INT_PREFIX_VAL    = "Int"
	JSON_UINT_PREFIX_VAL   = "UInt"
	JSON_FLOAT_PREFIX_VAL  = "Float"
	JSON_STRING_PREFIX_VAL = "String"
	JSON_BOOL_PREFIX_VAL   = "Bool"
)

const (
	EXPECTED_JSON_INT        = int8(-123)
	EXPECTED_JSON_UINT       = uint16(567)
	EXPECTED_JSON_FLOAT      = float32(67.8)
	EXPECTED_JSON_BOOL       = true
	EXPECTED_JSON_STRING     = "testdata"
	EXPECTED_JSON_INT_INT    = int16(5566)
	EXPECTED_JSON_INT_STRING = "intstring"
	EXPECTED_JSON_SLICE_SIZE = 3
	EXPECTED_JSON_SLICE_VAL0 = "one"
	EXPECTED_JSON_SLICE_VAL1 = "two"
	EXPECTED_JSON_SLICE_VAL2 = "free"
)

type IntJsonTestStruct struct {
	IntData    int16  `json:"IntIntData"`
	StringData string `json:"IntStrData"`
}

type JsonTestStruct struct {
	IntData    int8              `json:"${JSON_INT_PREFIX}Data"`
	UIntData   uint16            `json:"${JSON_UINT_PREFIX}Data"`
	FloatData  float32           `json:"${JSON_FLOAT_PREFIX}Data"`
	StringData string            `json:"${JSON_STRING_PREFIX}Data"`
	BoolData   bool              `json:"${JSON_BOOL_PREFIX}Data"`
	SliceData  []string          `json:"$.slice"`
	IntStruct  IntJsonTestStruct `json:"struct"`
}

const JSON_PROC_DATA = `{
		"root" : {
			"testcfg1" : {
				"IntData" : -123,
				"UIntData" : 567,
				"FloatData" : 67.8,
				"StringData" : "testdata",
				"BoolData" : true,
				"slice" : ["one","two","free"],
				"struct" : {
					"IntIntData" : 5566,
					"IntStrData" : "intstring"
				}
			}
		}
}
`

func TestCreateProcessor(t *testing.T) {
	var res any
	err := json.Unmarshal([]byte(JSON_DATA), &res)
	assert.NoError(t, err)

	// Case 1 correct create tag converter
	proc, err := NewJsonProcessor(res, TEST_PATH)
	assert.NoError(t, err)
	assert.NotNil(t, proc)

	// Case 2 json data is nil
	proc, err = NewJsonProcessor(nil, TEST_PATH)
	assert.Error(t, err)
	assert.Nil(t, proc)

	// Case 3 incorrect path
	proc, err = NewJsonProcessor(nil, INCORRECT_PATH)
	assert.Error(t, err)
	assert.Nil(t, proc)
}

func jsonProcessorVerifyResult(t *testing.T, res JsonTestStruct) {
	assert.Equal(t, EXPECTED_JSON_INT, res.IntData)
	assert.Equal(t, EXPECTED_JSON_UINT, res.UIntData)
	assert.Equal(t, EXPECTED_JSON_FLOAT, res.FloatData)
	assert.Equal(t, EXPECTED_JSON_STRING, res.StringData)
	assert.Equal(t, EXPECTED_JSON_BOOL, res.BoolData)
	assert.Equal(t, EXPECTED_JSON_SLICE_SIZE, len(res.SliceData))
	assert.Equal(t, EXPECTED_JSON_SLICE_VAL0, res.SliceData[0])
	assert.Equal(t, EXPECTED_JSON_SLICE_VAL1, res.SliceData[1])
	assert.Equal(t, EXPECTED_JSON_SLICE_VAL2, res.SliceData[2])
	assert.Equal(t, EXPECTED_JSON_INT_INT, res.IntStruct.IntData)
	assert.Equal(t, EXPECTED_JSON_INT_STRING, res.IntStruct.StringData)
}

func TestJsonConversion(t *testing.T) {
	var res any
	err := json.Unmarshal([]byte(JSON_PROC_DATA), &res)
	assert.NoError(t, err)
	jsonProcessor, err := NewJsonProcessor(res, "$.root.testcfg1")
	assert.NoError(t, err)
	assert.NotNil(t, jsonProcessor)
	jsonProcessor.SetDictionaryValue(JSON_INT_PREFIX_KEY, JSON_INT_PREFIX_VAL)
	jsonProcessor.SetDictionaryValue(JSON_UINT_PREFIX_KEY, JSON_UINT_PREFIX_VAL)
	jsonProcessor.SetDictionaryValue(JSON_FLOAT_PREFIX_KEY, JSON_FLOAT_PREFIX_VAL)
	jsonProcessor.SetDictionaryValue(JSON_STRING_PREFIX_KEY, JSON_STRING_PREFIX_VAL)
	jsonProcessor.SetDictionaryValue(JSON_BOOL_PREFIX_KEY, JSON_BOOL_PREFIX_VAL)
	//	defaultProcessor.SetDictionaryValue(INTERNAL_INT_PREFIX_KEY, INTERNAL_INT_PREFIX_VAL)
	//	defaultProcessor.SetDictionaryValue(INTERNAL_FLOAT_PREFIX_KEY, INTERNAL_FLOAT_PREFIX_VAL)
	//	defaultProcessor.SetDictionaryValue(INTERNAL_STRING_PREFIX_KEY, INTERNAL_STRING_PREFIX_VAL)
	testStruct := JsonTestStruct{}
	err = jsonProcessor.Process(&testStruct, nil)
	assert.NoError(t, err)
	jsonProcessorVerifyResult(t, testStruct)
}
