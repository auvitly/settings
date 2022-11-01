package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"settings/types"
)

func (c *Configurator) handleMap(handler *Handler) (err error) {

	var mapValue reflect.Value
	var mapType reflect.Type

	if err = c.checkUnsupportedMapTags(handler); err != nil {
		return err
	}

	switch c.options[types.ProcessingMode] {
	case types.ComplementMode:
		if !handler.reflectValue.IsNil() {
			mapValue = handler.reflectValue
		} else {
			mapType = reflect.MapOf(
				handler.reflectType.Key(),
				handler.reflectType.Elem(),
			)
			mapValue = reflect.MakeMapWithSize(mapType, 0)
		}
	case types.OverwritingMode:
		// Creating internal map
		mapType = reflect.MapOf(
			handler.reflectType.Key(),
			handler.reflectType.Elem(),
		)
		mapValue = reflect.MakeMapWithSize(mapType, 0)
	default:
		return ErrProcessing
	}

	for key, value := range handler.storage {

		// Creating internal storage for handler
		subStorage, ok := value.(map[string]interface{})
		if !ok {
			subStorage = make(map[string]interface{})
			// base value
			subStorage[key] = value
		}
		// Make internal handler for one record
		var fieldKey string
		for _, tag := range supportedTags {
			if _, ok = handler.fieldTags[tag]; ok {
				fieldKey = tag
			}
		}

		var subHandler *Handler

		subHandler = &Handler{
			name:           key,
			storage:        make(map[string]interface{}),
			reflectValue:   reflect.New(mapType.Elem()).Elem(),
			reflectType:    reflect.New(mapType.Elem()).Elem().Type(),
			structureField: reflect.StructField{},
			child:          make([]*Handler, 0),
			parent:         handler,
			fieldTags:      Tags{fieldKey: key},
			loadValues:     make(LoadValues),
			validator:      handler.validator,
		}

		switch handler.reflectValue.Type().Elem().Kind() {
		case reflect.Map, reflect.Pointer, reflect.Struct, reflect.Slice:
			if _, ok = subStorage[key].(map[string]interface{}); ok {
				subHandler.storage = subStorage
				break
			}
			if _, ok = subStorage[key].([]interface{}); ok {
				for idx, item := range subStorage[key].([]interface{}) {
					subHandler.storage[fmt.Sprintf("%d", idx)] = item
				}
				break
			}
			subHandler.storage = subStorage
		default:
			subHandler.storage = subStorage
		}

		if err = c.handle(subHandler); err != nil {
			return err
		}

		if handler.reflectValue.Type().Elem().Kind() == reflect.Pointer {
			mapValue.SetMapIndex(reflect.ValueOf(key), subHandler.reflectValue.Addr())
		} else {
			mapValue.SetMapIndex(reflect.ValueOf(key), subHandler.reflectValue)
		}
	}
	handler.reflectValue.Set(mapValue)
	return err

}

func (c *Configurator) checkUnsupportedMapTags(handler *Handler) (err error) {
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
