package internal

import (
	"github.com/pkg/errors"
	"net/url"
	"os"
	"reflect"
	"settings/types"
	"time"
)

func (c *Configurator) handleStruct(handler *Handler) (err error) {

	num := handler.reflectType.NumField()
	for i := 0; i < num; i++ {
		handler.parseStructField(i)
		switch handler.reflectValue.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Pointer, reflect.Map:
			kind := handler.reflectType.String()
			switch kind {
			case urlType:
				if err = c.handleUrl(handler); err != nil {
					return err
				}
			case timeType:
				if err = c.handleTime(handler); err != nil {
					return err
				}
			default:
				if err = c.handle(handler.child[i]); err != nil {
					return err
				}
			}
		default:
			err = ErrProcessing
		}
	}

	if handler.reflectType.String() == loggerType {
		if c.options[types.LoggerHook].(bool) {
			if configuration, ok := handler.reflectValue.Interface().(Logger); ok {
				c.configureLogger(configuration)
			}
		}
	}

	return
}

func (c *Configurator) handleUrl(h *Handler) (err error) {

	// Loading raw urlType
	for _, tag := range supportedTags {
		switch tag {
		case env:
			if result, ok := os.LookupEnv(h.fieldTags[tag]); ok {
				h.loadValues[tag] = result
			}
		case toml, yaml, xml, json:
			if result, ok := h.parent.storage[h.fieldTags[tag]]; ok {
				h.loadValues[tag] = result
			}
		default:
			if result, ok := h.fieldTags[tag]; ok {
				h.loadValues[tag] = result
			}
		}
	}

	// Getting raw urlType
	var rawUrl, tag string
	for _, tag = range supportedTags {
		if value, ok := h.loadValues[tag]; ok {
			rawUrl, ok = value.(string)
			if !ok {
				return errors.Wrapf(ErrBaseTypeNotMatch, "unsupported value: %v, define string for %s type field",
					h.obtainHandlerName(tag), urlType)
			}

			break
		}
	}
	if len(rawUrl) != 0 {
		url, err := url.Parse(rawUrl)
		if err != nil {
			return errors.Wrapf(err, "value: %v", h.obtainHandlerName(tag))
		}
		h.reflectValue.Set(reflect.ValueOf(url).Elem())
	}

	return nil

}

func (c *Configurator) handleTime(h *Handler) (err error) {

	// Loading raw urlType
	for _, tag := range supportedTags {
		switch tag {
		case env:
			if result, ok := os.LookupEnv(h.fieldTags[tag]); ok {
				h.loadValues[tag] = result
			}
		case toml, yaml, xml, json:
			if result, ok := h.parent.storage[h.fieldTags[tag]]; ok {
				h.loadValues[tag] = result
			}
		default:
			if result, ok := h.fieldTags[tag]; ok {
				h.loadValues[tag] = result
			}
		}
	}

	// Getting raw urlType
	var rawTime, tag string
	for _, tag = range supportedTags {
		if value, ok := h.loadValues[tag]; ok {
			rawTime, ok = value.(string)
			if !ok {
				return errors.Wrapf(ErrBaseTypeNotMatch, "unsupported value: %v, define string for %s type field",
					h.obtainHandlerName(tag), urlType)
			}

			break
		}
	}
	if len(rawTime) != 0 {
		opt, ok := c.options[types.TimeFormat].(string)
		if !ok {
			return ErrInvalidOptions
		}
		time, err := time.Parse(opt, rawTime)
		if err != nil {
			return errors.Wrapf(err, "value: %v", h.obtainHandlerName(tag))
		}
		h.reflectValue.Set(reflect.ValueOf(time))
	}
	return err
}

func (c *Configurator) handleLogger(h *Handler) (err error) {

	if err = c.handle(h); err != nil {
		return err
	}

	return nil
}
