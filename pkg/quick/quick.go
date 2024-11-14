package quick

import (
	"math"
	"math/rand/v2"
	"reflect"
	"unsafe"
)

func New[T any]() T {
	var value T
	concrete := reflect.ValueOf(&value).Elem()
	NewReflect(concrete)
	return value
}

func Update[T any](value *T) {
	concrete := reflect.ValueOf(value).Elem()
	NewReflect(concrete)
}

func NewReflect(value reflect.Value) {
	// TODO: Fix float generation to also include negative numbers.
	switch kind := value.Kind(); kind {
	case reflect.Bool:
		value.SetBool(rand.Int()&1 == 0)
	case reflect.Float32:
		value.SetFloat(float64(rand.Float32() * math.MaxFloat32))
	case reflect.Float64:
		value.SetFloat(rand.Float64() * math.MaxFloat64)
	case reflect.Complex64:
		value.SetComplex(complex(float64(rand.Float32()), float64(rand.Float32())))
	case reflect.Complex128:
		value.SetComplex(complex(rand.Float64(), rand.Float64()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(rand.Int64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value.SetUint(rand.Uint64())
	case reflect.Struct:
		n := value.NumField()
		for i := 0; i < n; i++ {
			field := value.Field(i)
			if field.CanSet() {
				NewReflect(field)
			} else {
				fieldPtr := unsafe.Pointer(field.UnsafeAddr())
				unsafeField := reflect.NewAt(field.Type(), fieldPtr).Elem()
				NewReflect(unsafeField)
			}
		}
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		NewReflect(value.Elem())
	}
}