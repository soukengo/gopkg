package copier

import "github.com/jinzhu/copier"


// Option sets copy options
type Option struct {
	// setting this value to true will ignore copying zero values of all the fields, including bools, as well as a
	// struct having all it's fields set to their zero values respectively (see IsZero() in reflect/value.go)
	IgnoreEmpty bool
	DeepCopy    bool
	Converters  []TypeConverter
}


func Copy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{Converters: converters})
}
func CopyWithOption(toValue interface{}, fromValue interface{}, option Option) error {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{DeepCopy: option.DeepCopy, IgnoreEmpty: option.IgnoreEmpty, Converters: converters})
}
