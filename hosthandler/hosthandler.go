package hosthandler

import (
	"net/http"
	"sync"

)

type HostHandler struct {
	HostHandlers    *sync.Map
	NotFoundHandler http.Handler
}

func New() *HostHandler {
	return &HostHandler{
		&sync.Map{},
		http.HandlerFunc(http.NotFound),
	}
}

func (d *HostHandler) Add(host string, h http.Handler) {
	d.HostHandlers.Store(host, h)
}

func (d *HostHandler) Remove(host string) {
	d.HostHandlers.Delete(host)
}

func (d *HostHandler) Get(host string) http.Handler {
	if h, ok := d.HostHandlers.Load(host); ok {
		return h.(http.Handler)
	}
	return nil
}

func (d *HostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := d.HostHandlers.Load(r.Host); !ok {
		d.NotFoundHandler.ServeHTTP(w, r)
	} else if h, ok := h.(http.Handler); ok {
		h.ServeHTTP(w, r)
	}
}

// Immutable creates a immutable copy of this router config.
func (d *HostHandler) Immutable() http.Handler {
	m := make(map[string]http.Handler)
	d.HostHandlers.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(http.Handler)
		return true
	})
	return NewImmutable(m, d.NotFoundHandler)
}

func NewImmutable(m map[string]http.Handler, h http.Handler) http.Handler {
	return &ImmutableHostHandler{m, h}
}

type ImmutableHostHandler struct {
	HostHandlers    map[string]http.Handler
	NotFoundHandler http.Handler
}

func (d *ImmutableHostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := d.HostHandlers[r.Host]; !ok {
		d.NotFoundHandler.ServeHTTP(w, r)
	} else {
		h.ServeHTTP(w, r)
	}
}
