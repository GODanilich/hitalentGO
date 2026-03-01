package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// - поле отсутствует в JSON -> Present=false
// - поле присутствует и null -> Present=true, Value=nil
// - поле присутствует и число -> Present=true, Value=&num
type OptionalInt64 struct {
	Present bool
	Value   *int64
}

func (o *OptionalInt64) UnmarshalJSON(b []byte) error {
	o.Present = true

	b = bytes.TrimSpace(b)
	if bytes.Equal(b, []byte("null")) {
		o.Value = nil
		return nil
	}

	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("must be int64 or null")
	}
	o.Value = &v
	return nil
}
