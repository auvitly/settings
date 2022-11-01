package internal

func (c *Configurator) handlePointer(handler *Handler) (err error) {
	err = handler.pointerFill()
	if err != nil {
		return
	}
	return c.handle(handler)
}
