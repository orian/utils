package reflect

import (
	"fmt"
	"testing"
)

type user struct {
	Name     string
	Age      int
	password string
}

func TestGetFieldNames(t *testing.T) {
	u := user{}
	cmpStrSl := func(a, b []string) string {
		if len(a) != len(b) {
			return fmt.Sprintf("len(a)!=len(b) (%d!=%d)", len(a), len(b))
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return fmt.Sprintf("a[%d]!=b[%d] (%q!=%q)", i, i, a[i], b[i])
			}
		}
		return ""
	}
	exp := []string{"Name", "Age", "password"}
	f, err := GetFieldNames(u)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if s := cmpStrSl(exp, f); s != "" {
		t.Errorf("literal field names don't match: %s", s)
	}

	f, err = GetFieldNames(u)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if s := cmpStrSl(exp, f); s != "" {
		t.Errorf("pointer field names don't match: %s", s)
	}

	sl := []string{"napis"}
	f, err = GetFieldNames(sl)
	if err == nil {
		t.Errorf("unexpected error: %s", err)
	}
}
