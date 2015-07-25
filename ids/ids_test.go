package ids

import (
	"testing"
)

func TestIds(t *testing.T) {
	str := Encode(1024)
	if str != "gBA=" {
		t.Errorf("want: gBA=, got: %s", str)
	}
	id, err := Decode(str)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if id != 1024 {
		t.Errorf("want: 1024, got: %d", id)
	}
}
