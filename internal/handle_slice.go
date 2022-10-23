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
	sliceTypeOf := reflect.SliceOf(handler.reflectValue.Type().Elem())
	sliceValueOf := reflect.MakeSlice(sliceTypeOf, 0, len(rawSlice))

	for index, value := range rawSlice {
		// Creating internal storage for handler
		subStorage, ok := value.(map[string]interface{})
		if !ok {
			return ErrHandle
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
			lv:             make(LoadValues),
			validator:      handler.validator,
		}
		if err = c.handle(subHandler); err != nil {
			return err
		}
		if handler.reflectValue.Type().Elem().Kind() == reflect.Pointer {
			sliceValueOf = reflect.Append(sliceValueOf, subHandler.reflectValue.Addr())
		} else {
			sliceValueOf = reflect.Append(sliceValueOf, subHandler.reflectValue)
		}
	}
	handler.reflectValue.Set(sliceValueOf)
	return
}
