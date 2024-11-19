package dynamictags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SIMPLE_STRING_NO_SUBST   = "Simple string"
	TEST_SUBST               = "TEST"
	SIMPLE_STRING_SUBST      = "Simple_${" + TEST_SUBST + "}_Tail"
	SIMPLE_SUBST_VALUE       = "VALUE"
	SIMPLE_SUBST_VALUE1      = "VAL345"
	EXPECTED_SIMPLE_NO_DICT  = "Simple__Tail"
	EXPECTED_SIMPLE_SUBST    = "Simple_" + SIMPLE_SUBST_VALUE + "_Tail"
	EXPECTED_SIMPLE_SUBST1   = "Simple_" + SIMPLE_SUBST_VALUE1 + "_Tail"
	LEVEL_2_KEY              = "LEVEL_1"
	LEVEL_2                  = "Level2"
	LEVEL_1_KEY              = "Level_" + LEVEL_2_KEY + "_LevelTail"
	LEVEL_1                  = "Level1"
	TWO_LEVEL_STRING         = "Simple_${" + LEVEL_1_KEY + "}_Tail"
	EXPECTED_TWO_LEVEL       = "Simple_" + LEVEL_1 + "_Tail"
	SIMPLE_NO_CLOSE_BRACE    = "Simple_${" + TEST_SUBST + "_Tail"
	SIMPLE_NO_CLOSE_BRACE1   = "}Simple_${" + TEST_SUBST + "_Tail"
	TWO_LEVEL_NO_CLOSE_BRACE = "Level_${${LEVEL2}_Tail"
	TWO_LEVEL_NO_TAIL        = "${${" + LEVEL_2_KEY + "}}"
)

func TestNoDictionary(t *testing.T) {
	//Case 1 no substitute,without dictionary
	res, err := ProcessString(SIMPLE_STRING_NO_SUBST, nil)
	assert.NoError(t, err)
	assert.Equal(t, SIMPLE_STRING_NO_SUBST, res)
	//Case 2 with one level substitute,no dictionary
	res, err = ProcessString(SIMPLE_STRING_SUBST, nil)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_SIMPLE_NO_DICT, res)
	//Case 3 with two levels, no dictionary
	res, err = ProcessString(TWO_LEVEL_STRING, nil)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_SIMPLE_NO_DICT, res)
}

func TestWithDictionary(t *testing.T) {
	dict := make(map[string]string)
	dict[TEST_SUBST] = SIMPLE_SUBST_VALUE
	dict[LEVEL_2_KEY] = LEVEL_2
	dict[LEVEL_1_KEY] = LEVEL_1
	//Case 1 no substitute
	res, err := ProcessString(SIMPLE_STRING_NO_SUBST, dict)
	assert.NoError(t, err)
	assert.Equal(t, SIMPLE_STRING_NO_SUBST, res)
	//Case 2 with one level substitute
	res, err = ProcessString(SIMPLE_STRING_SUBST, dict)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_SIMPLE_SUBST, res)
	//Case 3 with two levels
	res, err = ProcessString(TWO_LEVEL_STRING, dict)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_TWO_LEVEL, res)
	//Case 4 with two levels no tails
	dict[LEVEL_2] = LEVEL_1
	res, err = ProcessString(TWO_LEVEL_NO_TAIL, dict)
	assert.NoError(t, err)
	assert.Equal(t, LEVEL_1, res)
}

func TestEnvVariable(t *testing.T) {
	dict := make(map[string]string)
	// Case 1 dictionary value not set. Environment variable is set
	os.Setenv(TEST_SUBST, SIMPLE_SUBST_VALUE)
	res, err := ProcessString(SIMPLE_STRING_SUBST, dict)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_SIMPLE_SUBST, res)
	// Case 2 dictionary value is set. Environment value is set. Check dictionary value has priority
	dict[TEST_SUBST] = SIMPLE_SUBST_VALUE1
	res, err = ProcessString(SIMPLE_STRING_SUBST, dict)
	assert.NoError(t, err)
	assert.Equal(t, EXPECTED_SIMPLE_SUBST1, res)
}

func TestErrors(t *testing.T) {
	// Case 1 No close brace
	_, err := ProcessString(SIMPLE_NO_CLOSE_BRACE, nil)
	assert.Error(t, err)
	// Case 2 Incorrect position of close brace
	_, err = ProcessString(SIMPLE_NO_CLOSE_BRACE1, nil)
	assert.Error(t, err)
	// Case 3 No close brace on second level
	_, err = ProcessString(TWO_LEVEL_NO_CLOSE_BRACE, nil)
	assert.Error(t, err)
}
