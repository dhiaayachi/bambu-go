package jsonpatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func PatchValues(src []byte, iStructPointer interface{}) error {
	payloadMap := make(map[string]interface{})

	err := json.Unmarshal(src, &payloadMap)
	if err != nil {
		return err
	}

	structReflectValue, err := getReflectValueFromIStructPointer(iStructPointer)
	if err != nil {
		return err
	}

	err = traverseStructAndMergeStructFieldsWithPayload(structReflectValue, payloadMap)
	if err != nil {
		return err
	}

	return nil
}

func getReflectValueFromIStructPointer(iStructPointer interface{}) (ret reflect.Value, err error) {
	valueOfIStructPointer := reflect.ValueOf(iStructPointer)
	typeOfIStructPointer := reflect.TypeOf(iStructPointer)
	// Read Third Law here: https://blog.golang.org/laws-of-reflection
	// Pointer is needed as a patch operation would require mutation.
	// A direct call to Elem results in panic, thus the if statement block below.
	if k := valueOfIStructPointer.Kind(); k != reflect.Ptr {
		err = errors.New(fmt.Sprintf("%+v should be the pointer of struct.", typeOfIStructPointer))
		return
	}

	valueOfIStructPointerElem := valueOfIStructPointer.Elem()

	if k := valueOfIStructPointerElem.Type().Kind(); k != reflect.Struct {
		err = errors.New(fmt.Sprintf("%+v should be the struct type.", typeOfIStructPointer))
		return
	}

	// Below is a further (and definitive) check regarding settability in addition to checking whether it is a pointer earlier.
	if !valueOfIStructPointerElem.CanSet() {
		err = errors.New(fmt.Sprintf("%+v is unable to set the values.", typeOfIStructPointer))
		return
	}

	ret = valueOfIStructPointerElem

	return
}

func traverseStructAndMergeStructFieldsWithPayload(structReflectValue reflect.Value, payloadMap map[string]interface{}) error {
	for index := 0; index < structReflectValue.NumField(); index += 1 {
		structField := structReflectValue.Type().Field(index)
		structFieldJsonTag, err := getJsonStructTag(structField)
		if err != nil {
			return err
		}

		if iPayloadValue, ok := payloadMap[structFieldJsonTag]; ok {
			structFieldValue := structReflectValue.Field(index)
			err := mergePayloadToStructField(structFieldValue, iPayloadValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func mergePayloadToStructField(structFieldValue reflect.Value, iPayloadValue interface{}) (err error) {

	structFieldDataType := structFieldValue.Kind()

	switch structFieldDataType {
	case reflect.Struct:
		return mergePayloadToStructSF(structFieldValue, iPayloadValue)
	case reflect.Map:
		return mergePayloadToMapSF(structFieldValue, iPayloadValue)
	case reflect.Slice:
		return mergePayloadToSliceSF(structFieldValue, iPayloadValue)
	case reflect.Interface:
		return mergePayloadToInterfaceSF(structFieldValue, iPayloadValue)
	case reflect.Bool:
		return mergePayloadToBoolSF(structFieldValue, iPayloadValue)
	case reflect.String:
		return mergePayloadToStringSF(structFieldValue, iPayloadValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return mergePayloadToNumberSF(structFieldValue, iPayloadValue)
	default:
		err = errors.New(fmt.Sprintf("Unsupported type %+v.", structFieldDataType))
		return
	}

}

func mergePayloadToStructSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	payloadMap, ok := iPayloadValue.(map[string]interface{})
	if !ok {
		return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging.", structFieldDataType))
	}

	err = traverseStructAndMergeStructFieldsWithPayload(structFieldValue, payloadMap)
	if err != nil {
		return err
	}
	return nil
}

func mergePayloadToMapSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	structFieldType := structFieldValue.Type()

	var mapReflectValue reflect.Value
	if iPayloadValue == nil {
		mapReflectValue = reflect.MakeMap(structFieldType)
	} else {
		if payloadKind := reflect.ValueOf(iPayloadValue).Kind(); payloadKind != reflect.Map {
			return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldDataType))
		}

		mapReflectValue, err = getNewReflectValueMapWithPayloadValues(structFieldValue.Type(), iPayloadValue)
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(mapReflectValue)
	return nil
}

func mergePayloadToInterfaceSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	_, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	structFieldValue.Set(reflect.ValueOf(iPayloadValue))
	return nil
}

func mergePayloadToBoolSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	if iPayloadValue == nil {
		structFieldValue.SetBool(false)
	} else {
		if reflect.ValueOf(iPayloadValue).Kind() != reflect.Bool {
			return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldDataType))
		}
		structFieldValue.Set(reflect.ValueOf(iPayloadValue).Convert(structFieldValue.Type()))
	}

	return nil
}

func mergePayloadToStringSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	if iPayloadValue == nil {
		structFieldValue.SetString("")
	} else {
		if reflect.ValueOf(iPayloadValue).Kind() != reflect.String {
			return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldDataType))
		}
		structFieldValue.Set(reflect.ValueOf(iPayloadValue).Convert(structFieldValue.Type()))
	}

	return nil
}

func mergePayloadToNumberSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}

	if iPayloadValue == nil {
		structFieldValue.Set(reflect.ValueOf(0).Convert(structFieldValue.Type()))
	} else {
		if !isNumericValue(iPayloadValue) {
			return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldDataType))
		}
		structFieldValue.Set(reflect.ValueOf(iPayloadValue).Convert(structFieldValue.Type()))
	}

	return nil
}

func mergePayloadToSliceSF(structFieldValue reflect.Value, iPayloadValue interface{}) error {
	structFieldDataType, err := helperCheckSettabilityAndSFDataType(structFieldValue)
	if err != nil {
		return err
	}
	var sliceReflectValue reflect.Value
	if iPayloadValue == nil {
		emptyInterfaceSlice := make([]interface{}, 0)
		sliceReflectValue = makeNewSlice(structFieldDataType, emptyInterfaceSlice)
	} else {
		if payloadKind := reflect.ValueOf(iPayloadValue).Kind(); payloadKind != reflect.Slice {
			return errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldDataType))
		}

		sliceReflectValue, err = getNewReflectValueSliceWithPayloadValues(structFieldValue, iPayloadValue)
		if err != nil {
			return err
		}
	}
	structFieldValue.Set(sliceReflectValue)
	return nil
}

func helperCheckSettabilityAndSFDataType(structFieldValue reflect.Value) (structFieldDataType reflect.Type, err error) {
	if !structFieldValue.CanSet() {
		err = errors.New(fmt.Sprintf("CanSet() failed."))
		return
	}

	structFieldDataType = structFieldValue.Type()

	return
}

func getNewReflectValueMapWithPayloadValues(structFieldType reflect.Type, iPayloadValue interface{}) (reflect.Value, error) {
	newMap := reflect.MakeMap(structFieldType)
	payloadMap, ok := iPayloadValue.(map[string]interface{})
	if !ok {
		return newMap, errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
	}

	for k, v := range payloadMap {
		mapItemReflectKey, mapKeyReflectVal := getMapKeyAndValueReflectType(structFieldType, k, v)
		newMap.SetMapIndex(mapItemReflectKey, reflect.ValueOf(v).Convert(mapKeyReflectVal))
	}
	return newMap, nil
}

func getNewReflectValueSliceWithPayloadValues(structFieldValue reflect.Value, iPayloadValue interface{}) (sliceReflectValue reflect.Value, err error) {
	if !structFieldValue.CanSet() {
		err = errors.New(fmt.Sprintf("CanSet() failed."))
		return
	}

	interfaceSlice, skip, err := parseStructValueToInterfaceArray(iPayloadValue)
	if err != nil || skip {
		return
	}

	// Don't support mutiple data type in array
	err = checkMultipleDataTypeInPayloadArray(interfaceSlice)
	if err != nil {
		return
	}

	structFieldType := structFieldValue.Type()
	if structFieldType.Kind() == reflect.Invalid {
		err = errors.New(fmt.Sprintf("Invalid type! %+v", structFieldType))
		return
	}

	sliceReflectValue = makeNewSlice(structFieldType, interfaceSlice)
	k := structFieldType.Elem().Kind()
	switch k {
	case reflect.Struct:
		for index, ival := range interfaceSlice {
			nestedPayload, ok := ival.(map[string]interface{})
			if !ok {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}

			arrayItemAsStruct := reflect.Indirect(reflect.New(structFieldType.Elem()))
			err = traverseStructAndMergeStructFieldsWithPayload(arrayItemAsStruct, nestedPayload)
			if err != nil {
				return
			}
			sliceReflectValue.Index(index).Set(arrayItemAsStruct)
		}
	case reflect.Map:
		for index, ival := range interfaceSlice {
			payloadMap, ok := ival.(map[string]interface{})
			if !ok {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}
			arrayItemType := structFieldValue.Type().Elem()
			mapReflectVal, err := getNewReflectValueMapWithPayloadValues(arrayItemType, payloadMap)
			if err != nil {
				return sliceReflectValue, err
			}
			sliceReflectValue.Index(index).Set(mapReflectVal.Convert(structFieldType.Elem()))
		}
	case reflect.Interface:
		for index, ival := range interfaceSlice {
			nestedPayload, ok := ival.(interface{})
			if !ok {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}
			sliceReflectValue.Index(index).Set(reflect.ValueOf(nestedPayload).Convert(structFieldType.Elem()))
		}
	case reflect.Slice:
		for index, ival := range interfaceSlice {
			slicePayload, ok := ival.([]interface{})
			if !ok {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}

			arrayItemAsSlice := reflect.Indirect(reflect.New(structFieldType.Elem()))
			// WARN: recursion below.
			nestedSliceRefletValue, err := getNewReflectValueSliceWithPayloadValues(arrayItemAsSlice, slicePayload)
			if err != nil {
				return sliceReflectValue, err
			}
			sliceReflectValue.Index(index).Set(nestedSliceRefletValue)
		}
	case reflect.Bool, reflect.String:
		for index, ival := range interfaceSlice {
			if reflect.ValueOf(ival).Kind() != k {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}
			sliceReflectValue.Index(index).Set(reflect.ValueOf(ival).Convert(structFieldType.Elem()))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		for index, ival := range interfaceSlice {
			if !isNumericValue(ival) {
				err = errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", structFieldType))
				return
			}
			sliceReflectValue.Index(index).Set(reflect.ValueOf(ival).Convert(structFieldType.Elem()))
		}
	default:
		err = errors.New(fmt.Sprintf("Invalid payload data for %+v", structFieldType))
	}
	return
}

func isNumericValue(iPayloadValue interface{}) bool {
	switch reflect.ValueOf(iPayloadValue).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func getJsonStructTag(structField reflect.StructField) (string, error) {
	jsonTag := structField.Tag.Get("json")

	tags := strings.Split(jsonTag, ",")
	if len(tags) != 0 {
		jsonTag = tags[0]
	}

	if jsonTag == "" {
		return "", errors.New(fmt.Sprintf("Missing json tag in %+v struct field.", structField.Name))
	}
	return jsonTag, nil
}

func getMapKeyAndValueReflectType(structFieldType reflect.Type, k string, v interface{}) (mapKeyReflectValue reflect.Value, mapValueReflectType reflect.Type) {
	mapValueReflectType = reflect.ValueOf(v).Type()
	mapKeyReflectValue = reflect.ValueOf(k).Convert(structFieldType.Key())
	return
}

func makeNewSlice(sliceType reflect.Type, interfaces []interface{}) reflect.Value {
	return reflect.MakeSlice(sliceType, len(interfaces), cap(interfaces))
}

func parseStructValueToInterfaceArray(val interface{}) ([]interface{}, bool, error) {
	interfaceSlice, ok := val.([]interface{})
	if !ok {
		err := errors.New(fmt.Sprintf("Invalid payload data for %+v: incompatible for merging", interfaceSlice))
		return interfaceSlice, false, err
	}

	if len(interfaceSlice) == 0 {
		return interfaceSlice, true, nil
	}
	return interfaceSlice, false, nil
}

func checkMultipleDataTypeInPayloadArray(interfaceSlice []interface{}) error {
	var payloadArrayItemDataType reflect.Kind = reflect.Invalid
	for _, ival := range interfaceSlice {
		payloadArrayItemActualDataType := reflect.TypeOf(ival).Kind()
		if payloadArrayItemDataType != reflect.Invalid {
			if payloadArrayItemDataType != payloadArrayItemActualDataType {
				err := errors.New(fmt.Sprintf("Unable to support multiple data type in Array or Slice.", ival))
				return err
			}
		}
		payloadArrayItemDataType = payloadArrayItemActualDataType
	}

	return nil
}
