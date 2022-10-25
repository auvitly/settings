package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

func (c *Configurator) handleSlice(handler *Handler) (err error) {

	switch handler.reflectValue.Kind() {
	case reflect.Slice:
		var isByteSlice bool
		if isByteSlice, err = c.checkByteSlice(handler); err != nil {
			return err
		}
		if isByteSlice {
			return c.handleBytesSlice(handler)
		}
		return c.handleCommonSlice(handler)
	default:
		return ErrHandle
	}

}

func (c *Configurator) checkByteSlice(handler *Handler) (result bool, err error) {
	byteCase := handler.reflectValue.Type().Elem().Kind() == reflect.Uint8
	var currentType reflect.Type
	var currentTag string
	for tag, value := range handler.fieldTags {
		if str, ok := handler.parent.storage[value]; ok {
			currentType = reflect.TypeOf(str)
			currentTag = tag
		}
	}
	if currentType == nil {
		return false, nil
	}
	isString := currentType.Kind() == reflect.String
	if byteCase && !isString {
		return false, errors.Wrapf(ErrBaseTypeNotMatch, "can't convert to []%s field: %s ",
			handler.reflectValue.Type().Elem().Kind().String(),
			handler.obtainHandlerName(currentTag))
	}
	if byteCase && isString {
		return true, err
	}
	return false, nil
}

func (c *Configurator) handleBytesSlice(handler *Handler) (err error) {

	for _, value := range handler.parent.storage {
		str := value.(string)
		handler.reflectValue.SetBytes([]byte(str))
	}

	return err
}

func (c *Configurator) handleCommonSlice(handler *Handler) (err error) {

	var rawSlice []interface{}
	for _, value := range handler.fieldTags {
		if storage, ok := handler.parent.storage[value]; ok {
			rawSlice = storage.([]interface{})
			break
		}
	}

	var (
		sliceTypeOf  reflect.Type
		sliceValueOf reflect.Value
	)

	sliceTypeOf = reflect.SliceOf(handler.reflectValue.Type().Elem())
	switch handler.reflectValue.Len() {
	case 0:
		sliceValueOf = reflect.MakeSlice(sliceTypeOf, len(rawSlice), len(rawSlice))
	case len(rawSlice):
		sliceValueOf = handler.reflectValue
	default:
		length := handler.reflectValue.Len()
		for i := length; length < len(rawSlice); i++ {
			if handler.reflectValue.Type().Elem().Kind() == reflect.Pointer {
				handler.reflectValue = reflect.Append(handler.reflectValue, reflect.New(sliceTypeOf.Elem()).Elem().Addr())
			} else {
				handler.reflectValue = reflect.Append(handler.reflectValue, reflect.New(sliceTypeOf.Elem()).Elem())
			}
		}
		sliceValueOf = handler.reflectValue
	}

	length := handler.reflectValue.Len()
	for index := 0; index < length; index++ {

		// Creating internal storage for handler
		var subStorage = make(map[string]interface{})
		var ok bool
		if index < len(rawSlice) {
			subStorage, ok = rawSlice[index].(map[string]interface{})
			if !ok {
				return ErrHandle
			}
		}
		// Make internal handler for one record
		var fieldKey string
		for _, tag := range supportedTags {
			if _, ok = handler.fieldTags[tag]; ok {
				fieldKey = tag
			}
		}

		subHandler := &Handler{
			name:           fmt.Sprintf("%d", index),
			storage:        subStorage,
			reflectValue:   reflect.New(sliceTypeOf.Elem()).Elem(),
			reflectType:    reflect.New(sliceTypeOf.Elem()).Elem().Type(),
			structureField: reflect.StructField{},
			child:          make([]*Handler, 0),
			parent:         handler,
			fieldTags:      Tags{fieldKey: fmt.Sprintf("%d", index)},
			loadValues:     make(LoadValues),
			validator:      handler.validator,
		}
		if err = c.handle(subHandler); err != nil {
			return err
		}
		handler.child = append(handler.child, subHandler)
		if handler.reflectValue.Type().Elem().Kind() == reflect.Pointer {
			handler.reflectValue.Index(index).Set(subHandler.reflectValue.Addr())
		} else {
			handler.reflectValue.Index(index).Set(subHandler.reflectValue)
		}
	}
	//
	//for index, value := range rawSlice {
	//	// Creating internal storage for handler
	//	subStorage, ok := value.(map[string]interface{})
	//	if !ok {
	//		return ErrHandle
	//	}
	//	// Make internal handler for one record
	//	var fieldKey string
	//	for _, tag := range supportedTags {
	//		if _, ok = handler.fieldTags[tag]; ok {
	//			fieldKey = tag
	//		}
	//	}
	//	subHandler := &Handler{
	//		name:           fmt.Sprintf("%d", index),
	//		storage:        subStorage,
	//		reflectValue:   reflect.New(sliceTypeOf.Elem()).Elem(),
	//		reflectType:    reflect.New(sliceTypeOf.Elem()).Elem().Type(),
	//		structureField: reflect.StructField{},
	//		child:          make([]*Handler, 0),
	//		parent:         handler,
	//		fieldTags:      Tags{fieldKey: fmt.Sprintf("%d", index)},
	//		loadValues:     make(LoadValues),
	//		validator:      handler.validator,
	//	}
	//	if err = c.handle(subHandler); err != nil {
	//		return err
	//	}
	//	if handler.reflectValue.Type().Elem().Kind() == reflect.Pointer {
	//		handler.reflectValue.Index(index).Set(subHandler.reflectValue.Addr())
	//	} else {
	//		handler.reflectValue.Index(index).Set(subHandler.reflectValue)
	//	}
	//}
	handler.reflectValue.Set(sliceValueOf)
	return
}
