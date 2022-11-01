package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"settings/types"
)

func (c *Configurator) handleSlice(handler *Handler) (err error) {

	if err = c.checkUnsupportedSliceTags(handler); err != nil {
		return err
	}

	var isByteSlice bool
	if isByteSlice, err = c.checkByteSlice(handler); err != nil {
		return err
	}
	if isByteSlice {
		return c.handleBytesSlice(handler)
	}
	return c.handleCommonSlice(handler)

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

func (c *Configurator) checkUnsupportedSliceTags(handler *Handler) (err error) {
	var ok bool
	if _, ok = handler.fieldTags[env]; ok {
		return errors.Wrapf(ErrUnsupportedFieldTag, fmt.Sprintf("field %s -> %s; field_name: %s",
			env, handler.obtainHandlerName(env), handler.name))
	}
	if _, ok = handler.fieldTags[defaultValue]; ok {
		return errors.Wrapf(ErrUnsupportedFieldTag, fmt.Sprintf("field %s -> %s; field_name: %s",
			defaultValue, handler.obtainHandlerName(defaultValue), handler.name))
	}
	return
}

func (c *Configurator) handleBytesSlice(handler *Handler) (err error) {

	for tag, value := range handler.parent.storage {
		for _, key := range handler.fieldTags {
			if tag == key {
				str := value.(string)
				handler.reflectValue.SetBytes([]byte(str))
			}
		}
	}

	return err
}

func (c *Configurator) handleCommonSlice(handler *Handler) (err error) {

	var rawSlice []interface{}
	var rawStorage map[string]interface{}
	for _, value := range handler.fieldTags {
		if storage, ok := handler.parent.storage[value]; ok {
			switch storage.(type) {
			case []interface{}:
				rawSlice = storage.([]interface{})
			default:
				return errors.Wrap(ErrProcessing, "")
			}
			break
		}
	}

	var (
		sliceTypeOf  reflect.Type
		sliceValueOf reflect.Value
	)

	sliceTypeOf = reflect.SliceOf(handler.reflectValue.Type().Elem())
	switch c.options[types.ProcessingMode] {
	case types.ComplementMode:
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
	case types.OverwritingMode:
		sliceValueOf = reflect.MakeSlice(sliceTypeOf, len(rawSlice), len(rawSlice))
	}

	length := sliceValueOf.Len()
	for index := 0; index < length; index++ {

		// Creating internal storage for handler
		var subStorage map[string]interface{}
		var ok bool
		if index < len(rawSlice) {
			if rawStorage != nil {
				subStorage = rawStorage
			} else {
				subStorage, ok = rawSlice[index].(map[string]interface{})
				if !ok {
					subStorage = make(map[string]interface{})
					// base value
					subStorage[fmt.Sprintf("%d", index)] = rawSlice[index]
				}
			}
		}

		// Make internal handler for one record
		var fieldKey string
		for _, tag := range supportedTags {
			if _, ok = handler.fieldTags[tag]; ok {
				fieldKey = tag
				break
			}
		}

		// Make sub handler for item of slice
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
			sliceValueOf.Index(index).Set(subHandler.reflectValue.Addr())
		} else {
			sliceValueOf.Index(index).Set(subHandler.reflectValue)
		}
	}
	handler.reflectValue.Set(sliceValueOf)
	return
}
