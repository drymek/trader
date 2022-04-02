package cmpopts

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TransformJSON() cmp.Option {
	return cmpopts.AcyclicTransformer("TransformJSON", func(s []byte) interface{} {
		var v interface{}
		if err := json.Unmarshal(s, &v); err != nil {
			return s // use unparseable input as the output
		}
		return v
	})
}
