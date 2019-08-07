package headers

import (
	"net/http"
	"strings"
)

const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	AcceptLanguage = "Accept-Language"

	// Custom headers.
	XForwardedFor = "X-Forwarded-For" // https://en.wikipedia.org/wiki/X-Forwarded-For

	CacheControl    = "Cache-Control"
	ContentLength   = "Content-Length"
	ContentType     = "Content-Type"
	ContentEncoding = "Content-Encoding"
)

// TODO Concept and test.
//
//type HttpHeader string
//
//func (h HttpHeader) Get(header http.Header) string {
//	return header.Get(string(h))
//}
//
//func (h HttpHeader) RGet(r *http.Request) string {
//	return h.Get(r.Header)
//}

func GetClientIp(r *http.Request) string {
	if f := r.Header.Get(XForwardedFor); len(f) >= 7 { // 0.0.0.0 is 7 char length.
		ips := strings.SplitN(f, ",", 2)
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	ip := r.RemoteAddr
	if i := strings.Index(r.RemoteAddr, ":"); i > 0 {
		ip = ip[:i]
	}
	return ip
}
