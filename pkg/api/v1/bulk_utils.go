package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"

	v1types "github.com/astronomerio/event-api/pkg/types/v1"
	"github.com/imdario/mergo"
)

func gzipToBatch(b []byte) (batch v1types.Batch, err error) {
	gzData, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	defer gzData.Close()
	d, err := ioutil.ReadAll(gzData)
	if err != nil {
		return
	}
	err = json.Unmarshal(d, &batch)
	return
}

func mergeFields(dst, src interface{}) error {
	return mergo.Map(dst, src)
}
