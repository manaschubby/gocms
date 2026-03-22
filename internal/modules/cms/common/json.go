package common

import (
	"encoding/json"
	"strings"
)

func JsonNumberDecode(src []byte, dst any) error {
	decoder := json.NewDecoder(strings.NewReader(string(src)))
	decoder.UseNumber()
	return decoder.Decode(dst)
}
