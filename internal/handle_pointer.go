package internal

func (c *Configurator) handlePointer(handler *Handler) (err error) {

	if err = handler.pointerFill(); err != nil {
		return err
	}
	return c.handle(handler)
}
