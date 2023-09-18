package errefl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alexbathome/errefl/pkg/taglect"
)

func New2[T error](params ...interface{}) error {
	var err T
	var out = Err{}

	taglect.TaggedStruct[T, Err]{
		Tags: map[string]map[bool]taglect.OnTag[Err]{
			// If the struct field is tagged with "errefl" then this field
			// Can be used to provide a replacement value for the error message
			"errefl": {
				false:
				// False, the field is not anonymous (not embedded)
				func(sf reflect.StructField, ms *Err, s string) {
					if ms.replacements == nil {
						ms.replacements = make(map[string]string)
					}
					ms.replacements[s] = fmt.Sprintf("%v", params[sf.Index[0]-1])
				},
			},
			// If the struct field is tagged with "errtpl" then this field
			// Can be used to provide a template for the error message
			"errtpl": {
				true:
				// True, because the field is anonymous (embedded)
				func(f reflect.StructField, out *Err, value string) {
					for k, v := range out.replacements {
						value = strings.Replace(value, fmt.Sprintf("{%s}", k), v, 1)
					}
					out.Message = value
				},
			},
		},
		GoStruct: err,
	}.Process(&out)

	// Now that we have generated our embeddable error struct (errelf.Err)
	// We can set it back onto our original type for the consumer.
	vrefStructValue := reflect.ValueOf(&err).Elem()
	vrefStructFieldVal := vrefStructValue.FieldByName("Err")
	vrefStructFieldVal.Set(reflect.ValueOf(out))
	return err
}

func NewWrapped2[T error](inner error, params ...interface{}) error {
	outer := New2[T](params...)
	//TODO(alexb) - Fix this, it's broken, I can't for some reason
	// type cast to the Wrap interface, despite the method being available
	// on the embedded errelf.Err struct.
	if wrapper, ok := outer.(interface{ Wrap(error) }); ok {
		wrapper.Wrap(inner)
	}
	return outer
}
