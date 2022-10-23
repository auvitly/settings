package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

var validatorModule = validator.New()

type Handler struct {
	name           string
	storage        map[string]interface{}
	reflectValue   reflect.Value
	reflectType    reflect.Type
	structureField reflect.StructField
	child          []*Handler
	parent         *Handler
	fieldTags      Tags
	lv             LoadValues
	validator      *validator.Validate
}

func (c *Configurator) newRootHandler(value interface{}) (*Handler, error) {

	h := &Handler{
		name:         "root",                 // used
		storage:      c.config,               // used
		reflectValue: reflect.ValueOf(value), // used
		reflectType:  reflect.TypeOf(value),  // used
		child:        make([]*Handler, 0),    // used
		parent:       nil,                    // user
		fieldTags:    make(Tags),             // not used for root handler
		lv:           make(LoadValues),
		validator:    validatorModule, // used
	}

	// It is necessary to check the received value for the possibility of processing
	switch h.reflectValue.Kind() {
	case reflect.Pointer:
		if h.reflectValue.IsNil() {
			if !h.reflectValue.CanSet() {
				return nil, ErrNotAddressable
			}
			newValue := reflect.New(h.reflectValue.Type().Elem())
			h.reflectValue.Set(newValue)
		}
		h.reflectValue = h.reflectValue.Elem()
		h.reflectType = h.reflectType.Elem()
	case reflect.Struct:
		if h.parent == nil {
			return nil, ErrNotAddressable
		}
	default:
		return nil, ErrNotAStruct
	}

	// If pointer, but on simple type
	if h.reflectValue.Type().Kind() != reflect.Struct {
		return nil, ErrNotAStruct
	}

	return h, nil

}

func (h *Handler) pointerFill() error {
	switch h.reflectValue.Kind() {
	case reflect.Pointer:
		if h.reflectValue.IsNil() {
			if !h.reflectValue.CanSet() {
				return errors.Wrap(ErrNotAddressable, h.reflectType.String())
			}
			newValue := reflect.New(h.reflectValue.Type().Elem())
			h.reflectValue.Set(newValue)
		}
		h.reflectValue = h.reflectValue.Elem()
		h.reflectType = h.reflectType.Elem()
	}
	return nil
}

func (h *Handler) parseStructField(index int) {

	// Make Handler
	handler := &Handler{
		name:           h.reflectType.Field(index).Name,
		storage:        make(map[string]interface{}),
		reflectValue:   h.reflectValue.Field(index),
		reflectType:    h.reflectValue.Field(index).Type(),
		structureField: h.reflectType.Field(index),
		fieldTags:      make(Tags),
		lv:             make(LoadValues),
		parent:         h,
		validator:      validatorModule,
	}

	// Searching for tags from the list of allowed tags
	for _, tag := range supportedTags {
		if len(handler.structureField.Tag.Get(tag)) != 0 {
			handler.fieldTags[tag] = handler.structureField.Tag.Get(tag)
		}
		if result, ok := h.storage[handler.fieldTags[tag]].(map[string]interface{}); ok {
			handler.storage = result
		}
	}

	h.child = append(h.child, handler)

}

func (h *Handler) obtainEntireFieldName(tag string, i int) string {
	var tags []string
	for entity := h.child[i]; entity.parent != nil; entity = entity.parent {
		if value, ok := entity.fieldTags[tag]; ok {
			tags = append(tags, value)
		}
	}
	// swap rotation
	reverseSlice(tags)
	return strings.Join(tags, ".")
}

func (h *Handler) obtainHandlerName(tag string) string {
	var tags []string
	for entity := h; entity.parent != nil; entity = entity.parent {
		if value, ok := entity.fieldTags[tag]; ok {
			tags = append(tags, value)
		}
	}
	// swap rotation
	reverseSlice(tags)
	return strings.Join(tags, ".")
}
