package dynamictags

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
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
	tagpaths := make(map[string]string)
	err := processor.processStructure(t, v, "$", tagpaths, blackList)
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

func (processor DynamicTagProcessor) getUIntInterfaceValue(val any) (uint64, error) {
	valUInt, err := processor.convertUInt(val)
	if err == nil {
		return valUInt, nil
	}
	valInt, err := processor.convertInt(val)
	if err == nil {
		return uint64(valInt), nil
	}
	valFloat, err := processor.convertFloat(val)
	if err == nil {
		return uint64(valFloat), err
	}
	return 0, err
}

func (processor DynamicTagProcessor) getIntInterfaceValue(val any) (int64, error) {
	valInt, err := processor.convertInt(val)
	if err == nil {
		return valInt, nil
	}
	valUint, err := processor.convertUInt(val)
	if err == nil {
		return int64(valUint), nil
	}
	valFloat, err := processor.convertFloat(val)
	if err == nil {
		return int64(valFloat), err
	}
	return 0, err
}

func (processor DynamicTagProcessor) getFloatInterfaceValue(val any) (float64, error) {
	valFloat, err := processor.convertFloat(val)
	if err == nil {
		return valFloat, nil
	}
	valInt, err := processor.convertInt(val)
	if err == nil {
		return float64(valInt), nil
	}
	valUInt, err := processor.convertUInt(val)
	if err == nil {
		return float64(valUInt), err
	}
	return 0, err
}

func (processor DynamicTagProcessor) convertSliceInterface(src []interface{}) ([]string, error) {
	res := make([]string, 0, len(src))
	for _, srcData := range src {
		str, ok := srcData.(string)
		if !ok {
			return res, errors.New("incompatible slice elements type")
		}
		res = append(res, str)
	}
	return res, nil
}

func (processor DynamicTagProcessor) convertSliceString(src string) ([]string, error) {
	res := strings.Split(src, ",")
	return res, nil
}

func (processor DynamicTagProcessor) setInterfaceSimpleValue(t reflect.StructField, v reflect.Value, val any, path string) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valInt, err := processor.getIntInterfaceValue(val)
		if err != nil {
			return err
		}
		v.SetInt(valInt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valUInt, err := processor.getUIntInterfaceValue(val)
		if err != nil {
			return err
		}
		v.SetUint(valUInt)
	case reflect.Float32, reflect.Float64:
		valFloat, err := processor.getFloatInterfaceValue(val)
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
	case reflect.Slice:
		slice, ok := val.([]interface{})
		if ok {
			sliceVal, err := processor.convertSliceInterface(slice)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(sliceVal))
		}
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
	case reflect.Slice:
		sliceVal, err := processor.convertSliceString(val)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(sliceVal))
	default:
		return errors.New("unexpected key type")
	}
	return nil
}

func (processor DynamicTagProcessor) processSimpleType(t reflect.StructField, v reflect.Value, tagpaths map[string]string, path string) error {
	for _, converter := range processor.converters {
		tag := converter.GetTag()
		tagVal := t.Tag.Get(tag)
		if tagVal == "" {
			continue
		}
		res, err := ProcessString(tagVal, processor.dictionary)
		if err != nil {
			return err
		}
		tagPath, ok := tagpaths[tag]
		if !ok {
			tagPath = path
		}
		val, isSet, err := converter.GetSimpleValue(res, t, v, tagPath)
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

func (processor DynamicTagProcessor) fillTagsPath(t reflect.StructField, tagPaths map[string]string) error {
	for _, converter := range processor.converters {
		tag := converter.GetTag()
		tagVal := t.Tag.Get(tag)
		tagVal, err := ProcessString(tagVal, processor.dictionary)
		if err != nil {
			return err
		}
		currTagPath := tagVal
		if !strings.HasPrefix(tagVal, "$") {
			if tagVal == "" {
				tagVal = t.Name
			}
			var ok bool
			currTagPath, ok = tagPaths[tag]
			if !ok {
				currTagPath = "$"
			}
			currTagPath = currTagPath + "." + tagVal
		}
		tagPaths[tag] = currTagPath
	}
	return nil
}

func (processor DynamicTagProcessor) processStructure(t reflect.Type, v reflect.Value, path string, tagpaths map[string]string, blackList []string) error {
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
			continue
		}
		var err error = nil
		currPath := path + "." + fieldType.Name
		if blackList == nil || !slices.Contains(blackList, currPath) {
			if fieldValue.Kind() == reflect.Struct {
				err = processor.fillTagsPath(fieldType, tagpaths)
				if err != nil {
					return err
				}
				err = processor.processStructure(fieldType.Type, fieldValue, currPath, tagpaths, blackList)
			} else {
				err = processor.processSimpleType(fieldType, fieldValue, tagpaths, path)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
