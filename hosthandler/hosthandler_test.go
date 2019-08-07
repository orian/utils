package hosthandler

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"fmt"
	"math/rand"
)

func TestNew(t *testing.T) {
	h := New()
	var notFoundCalled bool
	h.NotFoundHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request){
		notFoundCalled = true
	})
	var exampleCalled bool
	h.Add("example.com", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		exampleCalled= true
	}))
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
	h.ServeHTTP(resp, req)
	assert.True(t, exampleCalled)

	req = httptest.NewRequest(http.MethodGet, "https://notfound.com", nil)
	h.ServeHTTP(resp, req)
	assert.True(t, notFoundCalled)
}

func BenchmarkHostHandler_ServeHTTP(b *testing.B) {
	h := New()
	var hosts []string
	for i:=0;i<15;i++ {
		host := fmt.Sprintf("mydomain%d.com", i)
		hosts = append(hosts, host)
		h.Add(host, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request){}))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		resp := httptest.NewRecorder()
		for pb.Next() {
			id := int(rand.Int31n(int32(len(hosts))))
			req := httptest.NewRequest(http.MethodGet, "https://"+hosts[id], nil)
			h.ServeHTTP(resp, req)
		}
	})
}

func BenchmarkImmutableHostHandler_ServeHTTP(b *testing.B) {
	h := New()
	var hosts []string
	for i:=0;i<3;i++ {
		host := fmt.Sprintf("mydomain%d.com", i)
		hosts = append(hosts, host)
		h.Add(host, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request){}))
	}

	im := h.Immutable()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		resp := httptest.NewRecorder()
		for pb.Next() {
			id := int(rand.Int31n(int32(len(hosts))))
			req := httptest.NewRequest(http.MethodGet, "https://"+hosts[id], nil)
			im.ServeHTTP(resp, req)
		}
	})
}