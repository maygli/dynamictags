package dynamictags

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	START_INCLUDE = "${"
	END_INCLUDE   = "}"
)

// Process string. Replace ${ENVIRONMENT_VARIABLE} by value from dictionary or by
// environment variable if no dictionary value found or replace by empty string
// if no dictionary and no environment variable found
// Parameters:
//   - str source string
//   - dictionary dictionary
//
// Returns:
//   - result string or error
func ProcessString(str string, dictionary map[string]string) (string, error) {
	stIndx := strings.Index(str, START_INCLUDE)
	if stIndx < 0 {
		return str, nil
	}
	endIndx := strings.LastIndex(str, END_INCLUDE)
	if endIndx < 0 || endIndx < stIndx {
		msg := fmt.Sprintf("incorrect tag structure. String '%s'. No closed brace", str)
		return "", errors.New(msg)
	}
	contentStr := str[stIndx+len(START_INCLUDE) : endIndx]
	content, err := ProcessString(contentStr, dictionary)
	if err != nil {
		return "", err
	}
	val, ok := dictionary[content]
	if !ok {
		val, ok = os.LookupEnv(content)
		if !ok {
			val = ""
		}
	}
	res := str[:stIndx] + val + str[endIndx+1:]
	return res, nil
}
