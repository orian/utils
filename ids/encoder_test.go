package ids

import (
	"bytes"
	"math"
	"testing"
)

func TestEncodeDecodeI32(t *testing.T) {
	be := &ByteEncoder{}
	bd := NewDecoder(&be.b)

	values := []int32{0, 1, 254, 255, 17, 123, 9538, 189231, math.MaxInt32, 1<<31 - 1}
	for i, v := range values {
		be.AppendInt32(v)

		g, ok := bd.ReadInt32()
		if !ok || g != v {
			t.Errorf("%d want: %d, got: %d", i, v, g)
		}
	}
}

func TestEncodeDecodeI64(t *testing.T) {
	be := &ByteEncoder{}
	bd := NewDecoder(&be.b)

	values := []int64{0, 17, 123, 9538, 189231, math.MaxInt64, 1<<63 - 1}
	for i, v := range values {
		be.AppendInt64(v)

		g, ok := bd.ReadInt64()
		if !ok || g != v {
			t.Errorf("%d want: %d, got: %d", i, v, g)
		}
	}
}

func TestEncodingSortOrder(t *testing.T) {
	be := &ByteEncoder{}
	values := [][]int32{
		[]int32{7},
		[]int32{123},
		[]int32{9538},
		[]int32{9538, 3},
		[]int32{189231},
		[]int32{math.MaxInt32},
	}
	var x [][]byte
	for _, v := range values {
		for _, v2 := range v {
			be.AppendInt32(v2)
		}
		a := be.Bytes()
		b := make([]byte, len(a))
		copy(b, a)
		x = append(x, b)
	}
	for i := 1; i < len(x); i++ {
		if bytes.Compare(x[i-1], x[i]) > 0 {
			t.Errorf("%d wrong order: %v %v", i, x[i-1], x[i])
		}
	}
}
