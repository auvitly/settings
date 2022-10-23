package internal

import (
	"reflect"
)

func (c *Configurator) handleMap(handler *Handler) (err error) {

	switch handler.reflectValue.Kind() {
	case reflect.Map:

		// Creating a internal map
		mapType := reflect.MapOf(
			handler.reflectType.Key(),
			handler.reflectType.Elem(),
		)
		mapValue := reflect.MakeMapWithSize(mapType, 0)

		for key, value := range handler.storage {
			switch handler.reflectValue.Type().Elem().Kind() {
			case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Struct:

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
					name:           key,
					storage:        subStorage,
					reflectValue:   reflect.New(mapType.Elem()).Elem(),
					reflectType:    reflect.New(mapType.Elem()).Elem().Type(),
					structureField: reflect.StructField{},
					child:          make([]*Handler, 0),
					parent:         handler,
					fieldTags:      Tags{fieldKey: key},
					lv:             make(LoadValues),
					validator:      handler.validator,
				}
				if err = c.handle(subHandler); err != nil {
					return err
				}
				mapValue.SetMapIndex(reflect.ValueOf(key), subHandler.reflectValue.Addr())
			default:
			}
		}
		handler.reflectValue.Set(mapValue)
	default:
		return ErrHandle
	}
	return

}
