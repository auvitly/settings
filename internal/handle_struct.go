package internal

import "reflect"

func (c *Configurator) handleStruct(handler *Handler) (err error) {
	for i := 0; i < handler.reflectType.NumField(); i++ {
		handler.parseStructField(i)
		switch handler.reflectValue.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Pointer, reflect.Map:
			if err = c.handle(handler.child[i]); err != nil {
				return err
			}
		default:
			err = c.handle(handler)
		}

	}
	return
}
