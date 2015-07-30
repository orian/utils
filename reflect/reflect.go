package reflect

import (
	"fmt"
	stdref "reflect"
)

var NotStructErr = fmt.Errorf("not a struct or pointer to struct")

func GetFieldNamesP(s interface{}) []string {
	a, err := GetFieldNames(s)
	if err != nil {
		panic(err)
	}
	return a
}

// Returns a field names using reflect package
func GetFieldNames(s interface{}) ([]string, error) {
	t := stdref.TypeOf(s)
	if t.Kind() == stdref.Ptr {
		t = t.Elem()
	}
	if kind := t.Kind(); kind != stdref.Struct {
		return nil, NotStructErr
	}
	var ret []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		ret = append(ret, field.Name)
	}
	return ret, nil
}
