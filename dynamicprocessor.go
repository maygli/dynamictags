package dynamictags

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

// Base struct for dynamic tag processors. Replace #{DICTIONATY_KEY} by value from
// processor dictionaty and ${ENVIRONMENT_VARIABLE} by environment variable
// For example we have processor (EnvProcessor) which process envoironment
// variables and struct which is defined in follow way:
//
//	type Data struct {
//	  EnvValue string `env:"${SERVER_NAME}_VALUE"`
//	}
//
// During processing:
//   - first name or environment variable will be calculated as value of
//     environment variable SERVER_NAME + '_VALUE' for example if
//     the environment variable 'SERVER_NAME' is set to 'TEST_SERVER' the result
//     name would be 'TEST_SERVER_VALUE'
//   - second value of 'EnvValue' will be set as value of environment variable
//     'TEST_SERVER_VALUE'.
type DynamicTagProcessor struct {
	dictionary map[string]string
	Tag        string
}

func (processor *DynamicTagProcessor) InitProcessor() {
	processor.dictionary = make(map[string]string)
}

func (processor *DynamicTagProcessor) SetDictionary(dict map[string]string) {
	processor.dictionary = dict
}

func (processor DynamicTagProcessor) GetDictionary() map[string]string {
	return processor.dictionary
}

func (processor *DynamicTagProcessor) SetDictionaryValue(key string, value string) {
	processor.dictionary[key] = value
}

func (processor *DynamicTagProcessor) RemoveDictionaryValue(key string) {
	delete(processor.dictionary, key)
}

func (processor DynamicTagProcessor) GetSimpleValue(tag string, t reflect.StructField, v reflect.Value, path string) (string, error) {
	return tag, nil
}

func (processor DynamicTagProcessor) processSimpleType(t reflect.StructField, v reflect.Value, path string) error {
	tag := t.Tag.Get(processor.Tag)
	if tag == "" {
		return nil
	}
	res, err := ProcessString(tag, processor.dictionary)
	if err != nil {
		return err
	}
	val, err := processor.GetSimpleValue(res, t, v, path)
	if err != nil {
		return err
	}
	switch t.Type.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(n) {
			return errors.New("int value overflow. Path: " + path)
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(n) {
			return errors.New("uint value overflow. Path: " + path)
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, v.Type().Bits())
		if err != nil {
			return err
		}
		if v.OverflowFloat(n) {
			return errors.New("float value overflow. Path: " + path)
		}
		v.SetFloat(n)
	case reflect.Bool:
		n, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		v.SetBool(n)
	default:
		return errors.New("unexpected key type")
	}
	fmt.Println("Processed value=", val)
	return nil
}

func (processor DynamicTagProcessor) processStructure(t reflect.Type, v reflect.Value, path string, blackList *[]string) error {
	var structValue reflect.Value
	var structType reflect.Type
	if v.Kind() == reflect.Pointer {
		structValue = v.Elem()
		structType = t.Elem()
	} else {
		structValue = v
		structType = t
	}
	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldType := structType.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s can not be changed", v.Elem().Type().Name())
		}
		var err error
		currPath := path + "." + fieldType.Name
		if blackList == nil || !slices.Contains(*blackList, currPath) {
			if fieldValue.Kind() == reflect.Struct {
				err = processor.processStructure(fieldType.Type, fieldValue, currPath, blackList)
			} else {
				err = processor.processSimpleType(fieldType, fieldValue, currPath)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor DynamicTagProcessor) Process(data any, blackList *[]string) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("pointer to structure is expected")
	}
	err := processor.processStructure(t, v, "", blackList)
	return err
}
