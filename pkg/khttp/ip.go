package khttp

import (
	"context"
	"net"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

func IpToContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return context.WithValue(ctx, contract.IpKey, realIP(r))
	}
}

func realIP(r *http.Request) string {
	var ip string
	var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}
