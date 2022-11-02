package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"settings/types"
	"strconv"
	"time"
)

// When processing simple values, several stages can be distinguished:
// 1. At the first stage, a type check is performed (structure, slice, etc.).
// If the type corresponds to the required one, then we proceed to the second stage
// 2. We load all available values into the loadValues handler variable

func (c *Configurator) handleBaseTypes(handler *Handler) (err error) {

	// Stage 1
	switch handler.reflectValue.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Struct:
		return ErrProcessing
	default:
		// Stage 2
		if err = handler.downloadTagValueBundles(); err != nil {
			return err
		}
		// Stage 3
		for key, value := range handler.fieldTags {
			storageValue := handler.loadValues[key]
			is := handler.isNeedToSetValue(key)
			if storageValue == nil || !is {
				continue
			}
			valueOfField := reflect.ValueOf(storageValue)
			switch kind := handler.reflectType.Kind(); kind {
			case reflect.Bool:
				if valueOfField.Kind() != handler.reflectValue.Kind() {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				handler.reflectValue.Set(valueOfField)
			case reflect.String:
				if valueOfField.Kind() != handler.reflectValue.Kind() {
					return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
				}
				handler.reflectValue.Set(valueOfField)
			case reflect.Float32, reflect.Float64:
				handler.reflectValue.SetFloat(valueOfField.Float())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch handler.reflectValue.Type().String() {
				case durationType:
					// Catch duration
					if err = handler.setDuration(storageValue); err != nil {
						return err
					}
				case syslogLevelType:
					lvl := storageValue.(types.SyslogLevel)
					handler.reflectValue.Set(reflect.ValueOf(lvl))
				default:
					// Processing with classic base types
					if valueOfField.Kind() != reflect.Float64 {
						return errBaserTypeNotMatch(value, handler.reflectType.Kind().String(), valueOfField.Kind().String())
					}
					if err = handler.checkingNumericalValueForLimit(storageValue.(float64), handler.reflectType.Kind(), key); err != nil {
						return err
					}
					f64 := storageValue.(float64)
					i64, err := strconv.ParseInt(fmt.Sprintf("%.f", f64), 10, 64)
					if err != nil {
						return err
					}
					handler.reflectValue.SetInt(i64)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				switch handler.reflectValue.Type().String() {
				case logrusLevelType:
					lvl := storageValue.(logrus.Level)
					handler.reflectValue.Set(reflect.ValueOf(lvl))
				default:
					if err = handler.checkingNumericalValueForLimit(storageValue.(float64), handler.reflectType.Kind(), key); err != nil {
						return err
					}
					f64 := storageValue.(float64)
					ui64, err := strconv.ParseUint(fmt.Sprintf("%.f", f64), 10, 64)
					if err != nil {
						return err
					}
					handler.reflectValue.SetUint(ui64)
				}
			default:
				return ErrUnsupportedFieldType
			}
		}
	}
	return

}

func (h *Handler) downloadTagValueBundles() error {
	// Crawl by handler tags
	for tag, field := range h.fieldTags {
		switch tag {
		// Search among env
		case env:
			if result, ok := os.LookupEnv(field); ok {
				h.loadValues[tag] = result
			}
		// Search among Viper values
		case toml, yaml, xml, json:
			if result, ok := h.parent.storage[field]; ok {
				switch kind := h.reflectValue.Type().String(); kind {
				case logrusLevelType:
					lvl, err := logrus.ParseLevel(result.(string))
					h.loadValues[tag] = lvl
					if err != nil {
						return err
					}
				case syslogLevelType:
					lvl, err := types.ParseSyslogPriority(result.(string))
					h.loadValues[tag] = lvl
					if err != nil {
						return err
					}
				default:
					h.loadValues[tag] = result
				}
			}
		// Detect default value
		case defaultValue:
			var err error
			if result, ok := h.fieldTags[tag]; ok {
				kind := h.reflectValue.Kind()
				switch kind {
				case reflect.String:
					switch kindT := h.reflectValue.Type().String(); kindT {
					default:
						h.loadValues[tag] = result
					}
				case reflect.Bool:
					h.loadValues[tag], err = strconv.ParseBool(result)
					if err != nil {
						return err
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					kindT := h.reflectValue.Type().String()
					switch kindT {
					case durationType:
						dur, err := time.ParseDuration(result)
						h.loadValues[tag] = float64(dur)
						if err != nil {
							return err
						}
					case logrusLevelType:
						lvl, err := logrus.ParseLevel(result)
						h.loadValues[tag] = lvl
						if err != nil {
							return err
						}
					case syslogLevelType:
						lvl, err := types.ParseSyslogPriority(result)
						h.loadValues[tag] = lvl
						if err != nil {
							return err
						}
					default:
						h.loadValues[tag], err = strconv.ParseFloat(result, 64)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (h *Handler) isNeedToSetValue(key string) bool {
	// Skip field with omit tag
	if _, ok := h.fieldTags[omit]; ok {
		return false
	}
	for _, tag := range supportedTags {
		// Else cases
		switch tag {
		case env:
			_, ok := h.loadValues[tag]
			if ok {
				if key == tag {
					return true
				}
				return false
			}
		case toml, yaml, xml, json:
			_, ok := h.loadValues[tag]
			if ok {
				if key == tag {
					return true
				}
				return false
			}
		default:
			_, ok := h.loadValues[tag]
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
					h.loadValues[tag] = result
				}
			case toml, yaml, xml, json:
				if result, ok := h.parent.storage[value]; ok {
					h.loadValues[tag] = result
				}
			case defaultValue:
				var err error
				if result, ok := h.fieldTags[tag]; ok {
					switch h.reflectValue.Kind() {
					case reflect.String:
						h.loadValues[tag] = result
					case reflect.Bool:
						h.loadValues[tag], err = strconv.ParseBool(result)
						if err != nil {
							return err
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						h.loadValues[tag], err = strconv.ParseInt(result, 10, 64)
						if err != nil {
							return err
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						h.loadValues[tag], err = strconv.ParseUint(result, 10, 64)
						if err != nil {
							return err
						}
					}

				}

			}
		}
	}
	return h.loadValues[key]
}

func (h *Handler) setDuration(storageValue interface{}) (err error) {
	if str, ok := storageValue.(string); ok {
		var dur time.Duration
		if dur, err = time.ParseDuration(str); err != nil {
			return err
		}
		h.reflectValue.SetInt(int64(dur))
		return nil
	}
	return nil
}

func (h *Handler) checkingNumericalValueForLimit(value float64, kind reflect.Kind, tag string) error {

	errFunc := func() error {
		return errors.Wrapf(ErrExceedingExpectValue, "Field name: %s | Field Type: %s | Found value: %.f",
			h.obtainHandlerName(tag),
			kind.String(),
			value,
		)
	}

	switch kind {
	case reflect.Int:
		if value > float64(maxInt) || value < float64(minInt) {
			return errFunc()
		}
	case reflect.Int8:
		if value > float64(maxInt8) || value < float64(minInt8) {
			return errFunc()
		}
	case reflect.Int16:
		if value > float64(maxInt16) || value < float64(minInt16) {
			return errFunc()
		}
	case reflect.Int32:
		if value > float64(maxInt32) || value < float64(minInt32) {
			return errFunc()
		}
	case reflect.Int64:
		if value > float64(maxInt64) || value < float64(minInt64) {
			return errFunc()
		}
	case reflect.Uint:
		if value > float64(maxUint) || value < float64(minUint) {
			return errFunc()
		}
	case reflect.Uint8:
		if value > float64(maxUint8) || value < float64(minUint8) {
			return errFunc()
		}
	case reflect.Uint16:
		if value > float64(maxUint16) || value < float64(minUint16) {
			return errFunc()
		}
	case reflect.Uint32:
		if value > float64(maxUint32) || value < float64(minUint32) {
			return errFunc()
		}
	case reflect.Uint64:
		if value > float64(maxUint64) || value < float64(minUint64) {
			return errFunc()
		}
	default:
		return ErrUnsupportedFieldType
	}

	return nil
}
