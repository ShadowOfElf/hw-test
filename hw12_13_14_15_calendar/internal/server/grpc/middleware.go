package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerLogRequestInterceptor(logg logger.LogInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		remoteAddr := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			remoteAddr = p.Addr.String()
		}

		md, _ := metadata.FromIncomingContext(ctx)
		userAgent := "unknown"
		if agent, exists := md["user-agent"]; exists && len(agent) > 0 {
			userAgent = agent[0]
		}

		logLine := fmt.Sprintf(
			"%s [%s] %s %s %s %d %d %s",
			remoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			"GRPC",
			info.FullMethod,
			"gRPC",
			getStatusCode(err),
			time.Since(start).Milliseconds(),
			userAgent,
		)
		logg.Info(logLine)
		return resp, err
	}
}

func getStatusCode(err error) int {
	if err != nil {
		return 500
	}
	return 200
}
