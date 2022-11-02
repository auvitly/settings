package internal

func (c *Configurator) handlePointer(handler *Handler) (err error) {

	if _, ok := handler.fieldTags[omit]; ok {
		return err
	}
	if err = handler.pointerFill(); err != nil {
		return err
	}
	return c.handle(handler)
}
