package ids

import (
	"bytes"
	"io"
)

type ByteEncoder struct {
	b, t bytes.Buffer
}

var toBase = func() []byte {
	x := make([]byte, 256)
	for i := 0; i < 255; i++ {
		x[i] = byte(i + 1)
	}
	x[255] = 0
	return x
}()

const encBase = 255
const marker = 0
const int32Size = 4
const int64Size = 8

func (b *ByteEncoder) AppendInt32(val int32) {
	if b.b.Len() > 0 {
		b.Mark()
	}
	b.t.Reset()
	for val > 0 {
		b.t.WriteByte(byte(val%encBase) + 1)
		val /= encBase
	}
	x := b.t.Bytes()
	s := len(x)
	for i := 0; i < int32Size-s; i++ {
		b.b.WriteByte(1)
	}
	for i := s - 1; i >= 0; i-- {
		b.b.WriteByte(x[i])
	}
}

func (b *ByteEncoder) AppendInt64(val int64) {
	if b.b.Len() > 0 {
		b.Mark()
	}
	b.t.Reset()
	for val > 0 {
		b.t.WriteByte(byte(val%encBase) + 1)
		val /= encBase
	}
	x := b.t.Bytes()
	s := len(x)
	for i := 0; i < int64Size-s; i++ {
		b.b.WriteByte(1)
	}
	for i := s - 1; i >= 0; i-- {
		b.b.WriteByte(x[i])
	}
}

func (b *ByteEncoder) Mark() {
	b.b.WriteByte(marker)
}

func (b *ByteEncoder) Bytes() []byte {
	x := b.b.Bytes()
	b.b.Reset()
	return x
}

type ByteDecoder struct {
	b io.ByteReader
}

func NewDecoder(b io.ByteReader) *ByteDecoder {
	return &ByteDecoder{b: b}
}

var powsOf255i32 = func() []int32 {
	x := make([]int32, int32Size)
	x[int32Size-1] = 1
	for i := int32Size - 2; i >= 0; i-- {
		x[i] = x[i+1] * 255
	}
	return x
}()

var powsOf255i64 = func() []int64 {
	x := make([]int64, int64Size)
	x[int64Size-1] = 1
	for i := int64Size - 2; i >= 0; i-- {
		x[i] = x[i+1] * 255
	}
	return x
}()

func (d *ByteDecoder) ReadInt32() (val int32, ok bool) {
	for i := 0; ; i++ {
		b, err := d.b.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, false
		}
		val += powsOf255i32[i] * int32(b-1)
	}
	return val, true
}

func (d *ByteDecoder) ReadInt64() (val int64, ok bool) {
	for i := 0; ; i++ {
		b, err := d.b.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, false
		}
		val += powsOf255i64[i] * int64(b-1)
	}
	return val, true
}
