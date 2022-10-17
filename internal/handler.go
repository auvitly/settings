package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	name      string
	rvalue    reflect.Value
	rtype     reflect.Type
	sfield    reflect.StructField
	child     []*Handler
	parent    *Handler
	tags      Tags
	lv        LoadValues
	validator *validator.Validate
}

func (h *Handler) validatePointer() error {
	switch h.rvalue.Kind() {
	case reflect.Pointer:
		if h.rvalue.IsNil() {
			if !h.rvalue.CanSet() {
				return ErrNotAddressable
			}
			newValue := reflect.New(h.rvalue.Type().Elem())
			h.rvalue.Set(newValue)
		}
		h.rvalue = h.rvalue.Elem()
		h.rtype = h.rtype.Elem()
		return nil
	case reflect.Struct:
		if h.rtype.NumField() == 0 {
			return ErrModelHasEmptyStruct
		}
		return nil
	default:
		return nil
	}
}

func (h *Handler) parseField(index int) {

	// Make Handler
	handler := &Handler{
		name:      h.rtype.Field(index).Name,
		rvalue:    h.rvalue.Field(index),
		rtype:     h.rvalue.Field(index).Type(),
		sfield:    h.rtype.Field(index),
		tags:      make(Tags),
		lv:        make(LoadValues),
		parent:    h,
		validator: h.validator,
	}

	// Searching for tags from the list of allowed tags
	for _, tag := range supportedTags {
		switch tag {
		case toml:
			if parent, ok := h.tags[tag]; ok && len(parent) == 0 {
				if len(handler.sfield.Tag.Get(tag)) != 0 {
					handler.tags[tag] = fmt.Sprintf("%s.%s", parent, handler.sfield.Tag.Get(tag))
				}
			} else {
				if len(handler.sfield.Tag.Get(tag)) != 0 {
					handler.tags[tag] = handler.sfield.Tag.Get(tag)
				}
			}
		default:
			if len(handler.sfield.Tag.Get(tag)) != 0 {
				handler.tags[tag] = handler.sfield.Tag.Get(tag)
			}
		}
	}
	h.child = append(h.child, handler)

}

func (h *Handler) ObtainEntireFieldName(tag string, i int) string {
	var tags []string
	for entity := h.child[i]; h.parent != nil && entity.tags != nil; entity = entity.parent {
		if value, ok := entity.tags[tag]; ok {
			tags = append(tags, value)
		}
	}
	reverse(tags)
	return strings.Join(tags, ".")
}
