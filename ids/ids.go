package ids

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

var RawURLEncoding = base64.URLEncoding.WithPadding(base64.NoPadding)

func Encode(id int64) string {
	b := make([]byte, 20)
	n := binary.PutVarint(b, id)
	return RawURLEncoding.EncodeToString(b[:n])
}

func Decode(str string) (int64, error) {
	data, err := RawURLEncoding.DecodeString(str)
	if err != nil {
		return 0, err
	}
	return binary.ReadVarint(bytes.NewReader(data))
}

