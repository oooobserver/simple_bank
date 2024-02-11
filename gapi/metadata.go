package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type MetaData struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetaData(ctx context.Context) *MetaData {
	meta := &MetaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}

		// This for extract useragent from the GRPC request
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			meta.UserAgent = userAgent[0]
		}

		if clientIp := md.Get(xForwardedForHeader); len(clientIp) > 0 {
			meta.ClientIp = clientIp[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIp = p.Addr.String()
	}

	return meta
}
