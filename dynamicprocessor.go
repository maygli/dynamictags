package dynamictags

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

// Base struct for dynamic tag processors.
type DynamicTagProcessor struct {
	dictionary map[string]string
	converters []TagConverterer
}

// Init dynamic processor
func (processor *DynamicTagProcessor) InitProcessor() {
	processor.dictionary = make(map[string]string)
}

// Set dictionary.
// Parameters:
//   - dict dictionary.
func (processor *DynamicTagProcessor) SetDictionary(dict map[string]string) {
	processor.dictionary = dict
}

// Returns dictionary.
// Returns:
//   - dictionary.
func (processor DynamicTagProcessor) GetDictionary() map[string]string {
	return processor.dictionary
}

// Set dictionary value.
// Parameters:
//   - key key
//   - value value
func (processor *DynamicTagProcessor) SetDictionaryValue(key string, value string) {
	processor.dictionary[key] = value
}

// Remove dictionary value.
// Parameters:
//   - key key
func (processor *DynamicTagProcessor) RemoveDictionaryValue(key string) {
	delete(processor.dictionary, key)
}

// Add tag converter
// Parameters:
//   - converter tag converter.
func (processor *DynamicTagProcessor) AddTagConverter(converter TagConverterer) {
	processor.converters = append(processor.converters, converter)
}

// Fill structure fields by tags
// Parameters:
//   - data pointer to filled structure.
//   - blackList black list of fields. This fields will be ignored during processing
//     Path to fields should be in the json path format. For example ($.InternalStructure.Field1)
//     Only simple paths is supported (without filters)
//
// Returns:
//   - error in case of error or nil
//
// For example we have processor (EnvProcessor) which process envoironment
// variables and struct which is defined in follow way:
//
//	type Data struct {
//	  EnvValue string `env:"${SERVER_NAME}_VALUE"`
//	}
//
// During processing:
//   - get value of 'SERVER_NAME' key from processor dictionary
//   - if the value is not exists get 'SERVER_NAME' environment variable
//   - if the environment variable is not defined use empty string
//   - replace 'SERVER_NAME' by the value. For example if 'SERVER_NAME' is set
//     to 'TEST' tag will be 'TEST_VALUE'
//   - call function processor.GetSimpleValue with tag 'TEST_VALUE' which calulate
//     'EnvValue' field value. Default implementation just returns tag value. This
//     behavouir can be changed in inhereted processor
func (processor DynamicTagProcessor) Process(data any, blackList []string) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("pointer to structure is expected")
	}
	err := processor.processStructure(t, v, "$", blackList)
	return err
}

func (processor DynamicTagProcessor) convertBool(val any) (bool, error) {
	res, ok := val.(bool)
	if ok {
		return res, nil
	}
	return false, errors.New("unsupported type")
}

func (processor DynamicTagProcessor) convertFloat(val any) (float64, error) {
	val32, ok := val.(float32)
	if ok {
		return float64(val32), nil
	}
	val64, ok := val.(float64)
	if ok {
		return val64, nil
	}
	return 0, errors.New("unsopported type")
}

func (processor DynamicTagProcessor) convertUInt(val any) (uint64, error) {
	valInt, ok := val.(uint)
	if ok {
		return uint64(valInt), nil
	}
	val8, ok := val.(uint8)
	if ok {
		return uint64(val8), nil
	}
	val16, ok := val.(uint16)
	if ok {
		return uint64(val16), nil
	}
	val32, ok := val.(uint32)
	if ok {
		return uint64(val32), nil
	}
	val64, ok := val.(uint64)
	if ok {
		return val64, nil
	}
	return 0, errors.New("unsopported type")
}

func (processor DynamicTagProcessor) convertInt(val any) (int64, error) {
	valInt, ok := val.(int)
	if ok {
		return int64(valInt), nil
	}
	val8, ok := val.(int8)
	if ok {
		return int64(val8), nil
	}
	val16, ok := val.(int16)
	if ok {
		return int64(val16), nil
	}
	val32, ok := val.(int32)
	if ok {
		return int64(val32), nil
	}
	val64, ok := val.(int64)
	if ok {
		return val64, nil
	}
	return 0, errors.New("unsopported type")
}

func (processor DynamicTagProcessor) setInterfaceSimpleValue(t reflect.StructField, v reflect.Value, val any, path string) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valInt, err := processor.convertInt(val)
		if err != nil {
			return err
		}
		v.SetInt(valInt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valUInt, err := processor.convertUInt(val)
		if err != nil {
			return err
		}
		v.SetUint(valUInt)
	case reflect.Float32, reflect.Float64:
		valFloat, err := processor.convertFloat(val)
		if err != nil {
			return err
		}
		v.SetFloat(valFloat)
	case reflect.Bool:
		valBool, err := processor.convertBool(val)
		if err != nil {
			return err
		}
		v.SetBool(valBool)
	}
	return nil
}

func (processor DynamicTagProcessor) setStringSimpleValue(t reflect.StructField, v reflect.Value, val string, path string) error {
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
	return nil
}

func (processor DynamicTagProcessor) processSimpleType(t reflect.StructField, v reflect.Value, path string) error {
	for _, converter := range processor.converters {
		tag := t.Tag.Get(converter.GetTag())
		if tag == "" {
			return nil
		}
		res, err := ProcessString(tag, processor.dictionary)
		if err != nil {
			return err
		}
		val, isSet, err := converter.GetSimpleValue(res, t, v, path)
		if err != nil {
			return err
		}
		if isSet {
			strVal, ok := val.(string)
			if ok {
				err = processor.setStringSimpleValue(t, v, strVal, path)
			} else {
				err = processor.setInterfaceSimpleValue(t, v, val, path)
			}
			return err
		}
	}
	return nil
}

func (processor DynamicTagProcessor) processStructure(t reflect.Type, v reflect.Value, path string, blackList []string) error {
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
		if blackList == nil || !slices.Contains(blackList, currPath) {
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
