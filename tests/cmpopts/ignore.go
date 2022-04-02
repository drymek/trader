package cmpopts

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func IgnoreFieldsContaining(fields []string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		return contains(fields, k)
	})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
