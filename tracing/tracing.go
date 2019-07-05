package tracing

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// DefaultTracingKeys is default tracing keys for istio
var DefaultTracingKeys = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}

// DefaultTracingKeys is default tracing keys for web server
var DefaultTracingKeysWeb = []string{
	"X-B3-Parentspanid",
	"X-B3-Sampled",
	"X-B3-Spanid",
	"X-B3-Traceid",
	"X-Request-Id",
	"X-B3-Flags",
	"X-Ot-Span-Context",
}

func mdPick(md metadata.MD, keys []string) metadata.MD {
	newMd := metadata.MD{}
	for _, key := range keys {
		v := md.Get(key)
		newMd.Set(key, v...)
	}
	return newMd
}

func include(arr []string, key string) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}

// Http2grpc pass tracing data from http header to downstream grpc metadata
func Http2grpc(ctx context.Context, tracingKeys []string, headers http.Header) context.Context {
	mp := map[string]string{}
	for _, key := range tracingKeys {
		v := headers.Get(key)
		if v != "" {
			mp[key] = v
		}
	}
	md := metadata.New(mp)
	return metadata.NewOutgoingContext(ctx, md)
}

// Grpc2http pass tracing data from grpc metadata to downstream http header
func Grpc2http(ctx context.Context, tracingKeys []string, dest http.Header) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	tracingMd := mdPick(md, tracingKeys)
	for k, v := range tracingMd {
		dest[k] = v
	}
}

// Grpc2Grpc pass tracing data from grpc metadata to downstream grpc metadata
func Grpc2Grpc(ctx context.Context, tracingKeys []string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	tracingMd := mdPick(md, tracingKeys)
	return metadata.NewOutgoingContext(ctx, tracingMd)
}

// Http2httpDest add tracing data to dest header
func Http2httpDest(tracingKeys []string, source, dest http.Header) {
	for k, v := range source {
		if include(tracingKeys, k) {
			dest[k] = v
		}
	}
}
