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

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}

	// NOTE: `md` is a map[string][]string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// This check is for the gateway server
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		// This check is for the plain GRPC server
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}

		// This check is for the gateway server
		if clients := md.Get(xForwardedForHeader); len(clients) > 0 {
			mtd.ClientIp = clients[0]
		}
	}

	// This check is for the plain GRPC server
	// NOTE: The IPAddress is stored in the context object of the request  not in the metadata
	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIp = p.Addr.String()
	}

	return mtd
}
