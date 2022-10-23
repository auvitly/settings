package internal

import (
	"github.com/pkg/errors"
	"os"
	"reflect"
	"time"
)

func (c *Configurator) handleBaseTypes(handler *Handler) (err error) {

	switch handler.reflectValue.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Struct:
		return ErrHandle
	default:
		for key, value := range handler.fieldTags {
			storageValue := handler.getValue(key, value)
			is := handler.isNeedToSetValue(key)
			if storageValue == nil || !is {
				continue
			}
			valueOfField := reflect.ValueOf(storageValue)
			switch handler.reflectType.Kind() {
			case reflect.String:
				if valueOfField.Kind() != handler.reflectValue.Kind() {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				handler.reflectValue.Set(valueOfField)
			case reflect.Float64:
				if valueOfField.Kind() != handler.reflectValue.Kind() {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				handler.reflectValue.Set(valueOfField)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// Catch duration setting
				isSet, err := handler.setDuration(storageValue)
				if err != nil {
					return err
				}
				if isSet {
					continue
				}
				// Processing with classic base types
				if valueOfField.Kind() != reflect.Float64 {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				handler.reflectValue.SetInt(int64(storageValue.(float64)))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if valueOfField.Kind() != reflect.Float64 {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				float64Value := storageValue.(float64)
				if float64Value < 0 {
					return errors.Wrapf(ErrBaseTypeNotMatch,
						"negative value: %s, define a positive value for %s type field",
						handler.obtainHandlerName(key),
						handler.reflectValue.Kind().String())
				}
				handler.reflectValue.SetUint(uint64(float64Value))
			default:
				return ErrUnsupportedFieldType
			}
		}
	}
	return

}

func (h *Handler) isNeedToSetValue(key string) bool {
	for _, tag := range supportedTags {
		switch tag {
		case env:
			_, ok := h.lv[tag]
			if ok {
				if key == tag {
					return true
				}
				return false
			}
		case toml, yaml, xml, json:
			_, ok := h.lv[tag]
			if ok {
				if key == tag {
					return true
				}
				return false
			}
		default:
			_, ok := h.lv[tag]
			if ok {
				if key == tag {
					return true
				}
				return false
			}
		}
	}
	return false
}

func (h *Handler) getValue(key, value string) interface{} {

	// load values
	for _, tag := range supportedTags {
		if key == tag {
			switch key {
			case env:
				if result, ok := os.LookupEnv(value); ok {
					h.lv[tag] = result
				}
			case toml, yaml, xml, json:
				if result, ok := h.parent.storage[value]; ok {
					h.lv[tag] = result
				}
			case defaultValue:
				if result, ok := h.fieldTags[tag]; ok {
					h.lv[tag] = result
				}
			}
		}
	}
	return h.lv[key]
}

func (h *Handler) setDuration(storageValue interface{}) (result bool, err error) {
	if str, ok := storageValue.(string); ok {
		var dur time.Duration
		if dur, err = time.ParseDuration(str); err != nil {
			return false, err
		}
		h.reflectValue.SetInt(int64(dur))
		return true, nil
	}
	return false, nil
}
