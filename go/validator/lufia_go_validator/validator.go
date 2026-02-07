package lufia_go_validator

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/lufia/go-validator"
)

func CollectErrors(err error) any {
	// pp.Printf("CollectErrors: type: %T, error:%v\n", err, err)

	rv := reflect.ValueOf(err)
	if rv.Kind() == reflect.Pointer && !rv.IsNil() {
		rv = rv.Elem()
	}

	// StructError[*any, any] ができないから、reflect で無理矢理やる。Unwrap してないけど。
	if rv.Kind() == reflect.Struct {
		errorsField := rv.FieldByName("Errors")
		if errorsField.IsValid() && errorsField.Type() == reflect.TypeFor[map[string]error]() {
			errs := errorsField.Interface().(map[string]error)
			m := map[string]any{}
			for k, e := range errs {
				m[k] = CollectErrors(e)
			}
			return m
		}
	}

	// この時点では SliceError はもう消えてるらしくて、このコードは動かない
	var sliceError *validator.SliceError[[]*Family, *Family]
	if errors.As(err, &sliceError) {
		m := map[string]any{}
		for _, k := range sliceError.Errors.Keys() {
			e, _ := sliceError.Errors.Get(k)
			m[strconv.Itoa(k)] = CollectErrors(e)
		}
		return m
	}

	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		var xs []any
		for _, e := range errs.Unwrap() {
			xs = append(xs, CollectErrors(e))
		}
		return xs
	}

	return err.Error()
}
