package taglect

import (
	"reflect"
)

// TaggedStruct is a helper struct to assist with the creation, and handling
// of complex Metaprogramming structs.
// It allows you to easily define what should happen when a tag is found. On a
// non meta struct.
type TaggedStruct[T, V any] struct {
	// Tags is a map of tags to functions that will be called when the tag is found.
	Tags TagsMap[V]
	// GoStruct is the struct that is literally written in the Go code
	// by the consumer.
	GoStruct T
}

// OnTag is a function that will be called when a tag is found.
type OnTag[V any] func(reflect.StructField, *V, string)

// TagsMap is a map of tags to functions that will be called when the tag is found.
// The boolean value is whether or not to skip the embedded field.
// It would look something like this:
//
//	{
//		"mytag": {
//			true: func(f reflect.StructField, out *MyStruct, value string) {
//				...
//			},
//		 	false: func(f reflect.StructField, out *MyStruct, value string) {
//				...
//			},
//		},
//	 }
type TagsMap[V any] map[string]map[bool]OnTag[V]

func (ts TaggedStruct[T, V]) rtype() reflect.Type {
	return reflect.TypeOf(ts.GoStruct)
}

// Fields iterates over the fields found from a struct using the reflect package.
// It returns a slice of reflect.StructField.
func (ts TaggedStruct[T, V]) Fields() []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < ts.rtype().NumField(); i++ {
		field := ts.rtype().Field(i)
		fields = append(fields, field)
	}
	return fields
}

// Process runs over the fields of the struct, and calls the functions defined
// by the consumer.
func (ts TaggedStruct[T, V]) Process(out *V) {
	for tag, tagMap := range ts.Tags {
		for anon, tagFn := range tagMap {
			for _, field := range ts.Fields() {
				if field.Anonymous && !anon {
					continue
				}
				tv := field.Tag.Get(tag)
				if tv != "" {
					tagFn(field, out, tv)
				}
			}
		}
	}
}
