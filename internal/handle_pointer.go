package internal

import "reflect"

func (c *Configurator) handlePointer(handler *Handler) (err error) {
	switch handler.reflectValue.Kind() {
	case reflect.Pointer:
		err = handler.pointerFill()
		if err != nil {
			return
		}
		return c.handle(handler)
	default:
		return ErrHandle
	}
}
