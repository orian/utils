package headers

import (
	"net/http"
	"testing"
)

func TestGetClientIp(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	r.RemoteAddr = "3.3.3.3"
	if err != nil {
		t.Error(err)
	}
	r.Header.Add(XForwardedFor, "10.0.0.10")
	r.Header.Add(XForwardedFor, "2.2.2.2")
	if ip := GetClientIp(r); ip != "10.0.0.10" {
		t.Errorf("want: 10.0.0.10, got: %s", ip)
	}
}
