package v1

import (
	"testing"
)

func Test_mergeFields(t *testing.T) {
	dst := map[string]interface{}{
		"val1": 100,
		"uniq": "string",
	}
	src := map[string]interface{}{
		"val1": 500,
		"val2": 100,
	}
	mergeFields(&dst, src)
	if dst["val2"] != 100 {
		t.Fatal("mergeFields didnt merge as expected")
	}
}
