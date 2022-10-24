package internal

import (
	"reflect"
)

func (c *Configurator) handle(handler *Handler) (err error) {

	kind := handler.reflectValue.Kind()
	switch kind {
	case reflect.Pointer:
		err = c.handlePointer(handler)
	case reflect.Struct:
		err = c.handleStruct(handler)
	case reflect.Map:
		err = c.handleMap(handler)
	case reflect.Slice:
		err = c.handleSlice(handler)
	default:
		err = c.handleBaseTypes(handler)
	}
	return err

}
