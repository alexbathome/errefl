package errefl

import (
	"fmt"
	"reflect"
	"strings"
)

type Err struct {
	Message  string
	InnerErr error

	replacements map[string]string
}

func (e Err) Unwrap() error {
	return e.InnerErr
}

func (e Err) Error() string {
	return e.Message
}

func (e *Err) Wrap(inner error) {
	e.InnerErr = inner
}

func NewWrapped[T error](inner error, details ...interface{}) error {
	err := New[T](details...).(T)
	embed := Err{
		Message:  err.Error(),
		InnerErr: inner,
	}
	setField(&err, "Err", embed)
	return err
}

func New[T error](items ...interface{}) error {
	var err T

	var messageFormatting map[string]string = make(map[string]string)

	// get our fields from the struct err
	e := reflect.TypeOf(err)

	// If we have specified an incorrect amount of items, we should just return a generic error, since the error
	// itself as not been defined correctly.
	// TODO(alexb) - This is SUPER inconvenient since all this reflection happens at runtime.
	// Find a way to do this at compile time.
	if len(items) != e.NumField()-1 {
		// We subtract one, since we assume that we have a errefl.Err field in the struct
		return fmt.Errorf("incorrect number of items passed to New. Expected %d, got %d", e.NumField()-1, len(items))
	}

	embeddedErrField, ok := e.FieldByName("Err")
	if !ok {
		return fmt.Errorf("could not find embedded Err field in struct")
	}
	embeddedStructVal := setupEmbedField(embeddedErrField)

	for i := 0; i < e.NumField(); i++ {
		// get the field
		f := e.Field(i)
		if f.Name == "Err" {
			// skip the embedded field,
			// TODO(alexb) - revisit this, as this isn't super resilient as a Type alias will break this.
			continue
		}

		// get the value
		value := items[i-1]
		setField(&err, f.Name, value)

		tag := f.Tag.Get("errefl")
		if tag != "" {
			messageFormatting[tag] = fmt.Sprintf("%v", value)
		}
	}

	for k, v := range messageFormatting {
		embeddedStructVal.Message = strings.Replace(embeddedStructVal.Message, fmt.Sprintf("{%s}", k), v, 1)
	}
	setField(&err, "Err", embeddedStructVal)

	return err
}

// setupEmbedField sets up the embedded errefl.Err field
func setupEmbedField(field reflect.StructField) Err {
	// get the tag
	var e Err
	tpl := field.Tag.Get("errtpl")
	if tpl != "" {
		if field.Type.Kind() == reflect.TypeOf(e).Kind() {
			e.Message = tpl
		}
	}

	return e
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func Catch[T error](v any) (T, bool) {
	o, ok := v.(T)
	return o, ok
}
