package copier

import (
	"github.com/jinzhu/copier"
)

var converters []copier.TypeConverter

type TypeConverter struct {
	SrcType interface{}
	DstType interface{}
	Fn      func(src interface{}) (interface{}, error)
}

func RegisterConverter(converter TypeConverter) {
	converters = append(converters,
		copier.TypeConverter{
			SrcType: converter.SrcType,
			DstType: converter.DstType,
			Fn:      converter.Fn,
		})
}
