package tracing

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// TracingKeys is default tracing keys for istio
var TracingKeys = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}

func mdPick(md metadata.MD, keys []string) metadata.MD {
	newMd := metadata.MD{}
	for _, key := range keys {
		v := md.Get(key)
		newMd.Set(key, v...)
	}
	return newMd
}

// Http2grpc pass tracing data from http header to downstream grpc metadata
func Http2grpc(ctx context.Context, headers http.Header) context.Context {
	mp := map[string]string{}
	for _, key := range TracingKeys {
		v := headers.Get(key)
		if v != "" {
			mp[key] = v
		}
	}
	md := metadata.New(mp)
	return metadata.NewOutgoingContext(ctx, md)
}

// Grpc2http pass tracing data from grpc metadata to downstream http header
func Grpc2http(ctx context.Context, originHeaders ...http.Header) http.Header {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return http.Header{}
	}
	tracingMd := mdPick(md, TracingKeys)
	tracingHeader := http.Header(tracingMd)
	for _, originHeader := range originHeaders {
		for k := range originHeader {
			tracingHeader.Set(k, originHeader.Get(k))
		}
	}
	return tracingHeader
}

// Grpc2Grpc pass tracing data from grpc metadata to downstream grpc metadata
func Grpc2Grpc(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	tracingMd := mdPick(md, TracingKeys)
	return metadata.NewOutgoingContext(ctx, tracingMd)
}

// Http2http pass tracing data from http header to downstream http header
func Http2http(headers http.Header, extHeaders ...http.Header) http.Header {
	mp := http.Header{}
	for _, key := range TracingKeys {
		v := headers.Get(key)
		if v != "" {
			mp[key] = []string{v}
		}
	}

	for _, extHeader := range extHeaders {
		for k, v := range extHeader {
			mp[k] = v
		}
	}

	return mp
}
