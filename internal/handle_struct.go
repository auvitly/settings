package internal

import (
	"github.com/pkg/errors"
	"net/url"
	"os"
	"reflect"
)

func (c *Configurator) handleStruct(handler *Handler) (err error) {
	for i := 0; i < handler.reflectType.NumField(); i++ {
		handler.parseStructField(i)
		switch handler.reflectValue.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Pointer, reflect.Map:
			if handler.reflectType.String() == urlType {
				return c.handleUrl(handler)
			}
			if err = c.handle(handler.child[i]); err != nil {
				return err
			}
		default:
			err = c.handle(handler)
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
				h.lv[tag] = result
			}
		case toml, yaml, xml, json:
			if result, ok := h.parent.storage[h.fieldTags[tag]]; ok {
				h.lv[tag] = result
			}
		default:
			if result, ok := h.fieldTags[tag]; ok {
				h.lv[tag] = result
			}
		}
	}

	// Getting raw urlType
	var rawUrl, tag string
	for _, tag = range supportedTags {
		if value, ok := h.lv[tag]; ok {
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
