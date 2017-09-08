package main

import (
	"github.com/upamune/grpc-error-sample/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	err = api.MashalError(ctx, err)
	return resp, err
}
