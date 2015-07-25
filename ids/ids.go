package ids

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

func Encode(id int64) string {
	b := make([]byte, 20)
	n := binary.PutVarint(b, id)
	return base64.StdEncoding.EncodeToString(b[:n])
}

func Decode(str string) (int64, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return 0, err
	}
	return binary.ReadVarint(bytes.NewReader(data))
}

