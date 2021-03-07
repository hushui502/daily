package converter

import (
	"bytes"
	"encoding/json"
)

// Convert bytes to buffer helper
func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}

	return buf, nil
}
