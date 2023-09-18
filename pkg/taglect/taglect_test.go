package taglect_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/alexbathome/errefl/pkg/taglect"
)

type MyErrorStruct struct {
	Err      error  `errtpl:"*{endpoint}* is not responding on: {protocol}"`
	Endpoint string `errefl:"endpoint"`
}

type MetaStruct struct {
	Message error

	replacements map[string]string
}

func CreateMetastruct[T any](t T) *MetaStruct {
	ttags := taglect.TaggedStruct[T, MetaStruct]{
		Tags: map[string]map[bool]taglect.OnTag[MetaStruct]{
			"errefl": {
				false: func(sf reflect.StructField, ms *MetaStruct, s string) {
					if ms.replacements == nil {
						ms.replacements = make(map[string]string)
					}
					ms.replacements[s] = fmt.Sprintf("%v", reflect.ValueOf(t).FieldByName(sf.Name))
				},
			},
			"errtpl": {
				false: func(f reflect.StructField, out *MetaStruct, value string) {
					for k, v := range out.replacements {
						value = strings.Replace(value, fmt.Sprintf("{%s}", k), v, 1)
					}
					out.Message = fmt.Errorf(value)
				},
			},
		},
		GoStruct: t,
	}
	var out = new(MetaStruct)
	ttags.Process(out)
	return out
}

func TestTagLect(t *testing.T) {
	err := CreateMetastruct[MyErrorStruct](MyErrorStruct{
		Endpoint: "https://example.com",
	})
	fmt.Println(err.Message)
}
