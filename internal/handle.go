package internal

import (
	"reflect"
)

func (c *Configurator) handle(handler *Handler) (err error) {

	kind := handler.reflectValue.Kind()
	switch kind {
	case reflect.Pointer:
		err = c.handlePointer(handler)
	case reflect.Struct:
		err = c.handleStruct(handler)
	case reflect.Map:
		err = c.handleMap(handler)
	case reflect.Slice:
		err = c.handleSlice(handler)
	default:
		err = c.handleBaseTypes(handler)
	}
	return err

}

func (c *Configurator) loadValues(handler *Handler, i int) (err error) {
	for tag, _ := range handler.child[i].fieldTags {
		switch tag {
		case env:
			//if result, ok := os.LookupEnv(value); ok {
			//	handler.child[i].loadValues[tag] = result
			//}
		case toml, yaml, json, xml:
			//scan := handler.obtainEntireFieldName(tag, i)
			//if result := c.viper.Get(scan); result != nil {
			//	handler.child[i].loadValues[tag] = result
			//}
		case defaultValue:
			//handler.child[i].loadValues[defaultValue] = handler.child[i].fieldTags[defaultValue]
		default:
			continue
		}
	}
	return nil
}

func (c *Configurator) settingValues(handler *Handler, i int) (err error) {

	//// Examining all support tags
	//for _, tag := range supportedTags {
	//
	//	// Does the tag contain inside loaded values
	//	value, isEnabled := handler.child[i].loadValues[tag]
	//	if !isEnabled {
	//		continue
	//	}
	//
	//	kind := handler.child[i].reflectValue.Kind()
	//	switch kind {
	//	case reflect.String:
	//		if result, ok := value.(string); ok {
	//			handler.child[i].reflectValue.SetString(result)
	//		}
	//		return
	//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//		if result, ok := value.(int64); ok {
	//			handler.child[i].reflectValue.SetInt(result)
	//		}
	//		return
	//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//		if result, ok := value.(uint64); ok {
	//			handler.child[i].reflectValue.SetUint(result)
	//		}
	//		return
	//	case reflect.Bool:
	//		if result, ok := value.(bool); ok {
	//			handler.child[i].reflectValue.SetBool(result)
	//		}
	//		return
	//	case reflect.Float32, reflect.Float64:
	//		if result, ok := value.(float64); ok {
	//			handler.child[i].reflectValue.SetFloat(result)
	//		}
	//		return
	//	default:
	//		return ErrUnsupportedFieldType
	//	}
	//
	//}
	//
	return nil
}
