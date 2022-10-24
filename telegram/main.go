package telegram

import (
	"bytes"
	"encoding/json"
	"time"
)

var (
	DefaultApiUrl = `https://api.telegram.org`
	DefaultRequestTimeout = time.Second * 3
	DefaultRequestRetries = 3
)

func JsonEncode(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encode := json.NewEncoder(buffer)
	if err := encode.Encode(v); err != nil { return nil, err }
	return buffer.Bytes(), nil
}
