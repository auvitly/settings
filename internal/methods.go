package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"os"
	"reflect"
)

func (c *Configurator) handle(config any) (err error) {

	// A function calls itself repeatedly (recursion).
	// Since a handler can be passed as an argument, the conversion can be omitted.
	handler, ok := config.(*Handler)
	if !ok {
		// Main handler without partner handler
		handler = &Handler{
			rvalue:    reflect.ValueOf(config),
			rtype:     reflect.TypeOf(config),
			child:     make([]*Handler, 0),
			parent:    nil,
			lv:        make(LoadValues),
			validator: validator.New(),
		}
	}

	// Jump from pointer to value
	if err = handler.validatePointer(); err != nil {
		return
	}

	// Let's traverse the structure
	for i := 0; i < handler.rtype.NumField(); i++ {
		switch handler.parseField(i); handler.child[i].rvalue.Kind() {
		case reflect.Struct:
			if err = c.handle(handler.child[i]); err != nil {
				return err
			}
			continue
		case reflect.Pointer:
			if err = c.processPointer(handler, i); err != nil {
				return err
			}
			if err = c.loadValues(handler, i); err != nil {
				return err
			}
			if err = c.settingValues(handler, i); err != nil {
				return err
			}
			continue
		default:
			// Setting the value from the environment variable
			if err = c.loadValues(handler, i); err != nil {
				return err
			}
			if err = c.settingValues(handler, i); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func (c *Configurator) loadValues(handler *Handler, i int) (err error) {
	for tag, value := range handler.child[i].tags {
		switch tag {
		case env:
			if result, ok := os.LookupEnv(value); ok {
				handler.child[i].lv[tag] = result
			}
		case toml, yaml, json:
			scan := handler.ObtainEntireFieldName(tag, i)
			if result := c.viper.Get(scan); result != nil {
				handler.child[i].lv[tag] = result
			}
		case defaultValue:
			handler.child[i].lv[defaultValue] = handler.child[i].tags[defaultValue]
		default:
			continue
		}
	}
	return nil
}

func (c *Configurator) settingValues(handler *Handler, i int) (err error) {

	// Examining all support tags
	for _, tag := range supportedTags {

		// Does the tag contain inside loaded values
		value, isEnabled := handler.child[i].lv[tag]
		if !isEnabled {
			continue
		}

		switch handler.child[i].rvalue.Kind() {
		case reflect.String:
			if result, ok := value.(string); ok {
				handler.child[i].rvalue.SetString(result)
			}
			return
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if result, ok := value.(int64); ok {
				handler.child[i].rvalue.SetInt(result)
			}
			return
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if result, ok := value.(uint64); ok {
				handler.child[i].rvalue.SetUint(result)
			}
			return
		case reflect.Bool:
			if result, ok := value.(bool); ok {
				handler.child[i].rvalue.SetBool(result)
			}
			return
		case reflect.Float32, reflect.Float64:
			if result, ok := value.(float64); ok {
				handler.child[i].rvalue.SetFloat(result)
			}
			return
		case reflect.Map:
			err = c.settingMap(tag, handler, i)
		case reflect.Slice:
			switch handler.child[i].rvalue.Type().Elem().Kind() {
			case reflect.String:
				result := c.viper.GetStringSlice(handler.ObtainEntireFieldName(tag, i))
				sliceTypeOf := reflect.SliceOf(handler.child[i].rvalue.Type().Elem())
				sliceValueOf := reflect.MakeSlice(sliceTypeOf, 0, len(result))
				for j := 0; j < len(result); j++ {
					switch handler.child[i].rvalue.Type().Elem().Kind() {
					case reflect.String:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(result[j]))
					default:
						return errors.Wrapf(ErrUnsupportedFieldType, "%v", handler.child[i].rvalue.Type().Elem().Kind())
					}
				}
				handler.child[i].rvalue.Set(sliceValueOf)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result := c.viper.GetIntSlice(handler.ObtainEntireFieldName(tag, i))
				sliceTypeOf := reflect.SliceOf(handler.child[i].rvalue.Type().Elem())
				sliceValueOf := reflect.MakeSlice(sliceTypeOf, 0, len(result))
				for j := 0; j < len(result); j++ {
					switch handler.child[i].rvalue.Type().Elem().Kind() {
					case reflect.Int:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(result[j]))
					case reflect.Int8:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(int8(result[j])))
					case reflect.Int16:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(int16(result[j])))
					case reflect.Int32:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(int32(result[j])))
					case reflect.Int64:
						sliceValueOf = reflect.Append(sliceValueOf, reflect.ValueOf(int64(result[j])))
					default:
						return errors.Wrapf(ErrUnsupportedFieldType, "%v", handler.child[i].rvalue.Type().Elem().Kind())
					}

				}
				handler.child[i].rvalue.Set(sliceValueOf)
			default:
				return ErrUnsupportedFieldType
			}
		default:
			return ErrUnsupportedFieldType
		}

	}

	return nil
}

func (c *Configurator) settingMap(tag string, handler *Handler, index int) error {

	result := c.viper.GetStringMap(handler.ObtainEntireFieldName(tag, index))

	if result != nil {

		var (
			viperValue = reflect.ValueOf(result)

			handlerElemType    = handler.child[index].rvalue.Type().Elem().Kind()
			viperValueElemType = viperValue.Type().Elem().Kind()
		)

		switch handlerElemType {
		// Equal element type
		case viperValueElemType:
			handler.child[index].rvalue.Set(viperValue)
			// Viper Float64 -> Int
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			// Creating a new collection
			mapType := reflect.MapOf(
				handler.child[index].rtype.Key(),
				handler.child[index].rtype.Elem(),
			)
			mapValue := reflect.MakeMapWithSize(mapType, 0)
			// Iteration
			iter := viperValue.MapRange()
			for iter.Next() {
				f, _ := iter.Value().Interface().(float64)
				var v reflect.Value
				switch handler.child[index].rtype.Elem().Kind() {
				case reflect.Int:
					v = reflect.ValueOf(int(f))
				case reflect.Int8:
					v = reflect.ValueOf(int8(f))
				case reflect.Int16:
					v = reflect.ValueOf(int16(f))
				case reflect.Int32:
					v = reflect.ValueOf(int32(f))
				case reflect.Int64:
					v = reflect.ValueOf(int64(f))
				case reflect.Uint:
					v = reflect.ValueOf(uint(f))
				case reflect.Uint8:
					v = reflect.ValueOf(uint8(f))
				case reflect.Uint16:
					v = reflect.ValueOf(uint16(f))
				case reflect.Uint32:
					v = reflect.ValueOf(uint32(f))
				case reflect.Uint64:
					v = reflect.ValueOf(uint64(f))
				case reflect.Float32:
					v = reflect.ValueOf(float32(f))
				case reflect.Float64:
					v = reflect.ValueOf(f)
				}
				mapValue.SetMapIndex(iter.Key(), v)
			}
			// Setting map
			handler.child[index].rvalue.Set(mapValue)
		default:
			return ErrUnsupportedFieldType
		}
	}
	return nil
}

func (c *Configurator) processPointer(handler *Handler, i int) (err error) {
	switch handler.child[i].rvalue.Type().Elem().Kind() {
	case reflect.Struct:
		if err = c.handle(handler.child[i]); err != nil {
			return err
		}
		break
	default:
		// Create a new variable if necessary
		if handler.child[i].rvalue.IsNil() {
			if !handler.child[i].rvalue.CanSet() {
				return ErrNotAddressable
			}
			newValue := reflect.New(handler.child[i].rvalue.Type().Elem())
			handler.child[i].rvalue.Set(newValue)
		}
		handler.child[i].rvalue = handler.child[i].rvalue.Elem()
		handler.child[i].rtype = handler.child[i].rtype.Elem()
		break
	}
	return nil
}
