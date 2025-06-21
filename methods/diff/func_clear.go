package diff

import "reflect"

// // // // // // // // // // // // // // // // // //

func clearObj(data any) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		panic("clearStruct requires a pointer")
	}

	elem := val.Elem()
	_clearValue(elem)
}

func _clearValue(value reflect.Value) {
	switch value.Kind() {

	case reflect.Ptr:
		if !value.IsNil() {
			_clearValue(value.Elem())
		}

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanSet() {
				_clearValue(field)
			}
		}

	case reflect.Interface:
		if !value.IsNil() {
			elem := value.Elem()
			_clearValue(elem)
			// value.Set(reflect.Zero(value.Type()))
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			bufVal := value.MapIndex(key)

			_, ok := bufVal.Interface().(string)
			if ok {
				value.SetMapIndex(key, reflect.Zero(reflect.TypeOf("")))
			} else {
				_clearValue(bufVal)
			}
		}

	default:
		if !value.IsValid() || !value.CanSet() {
			return
		}
		value.Set(reflect.Zero(value.Type()))
	}
}
