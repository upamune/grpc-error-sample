package api

import (
	"fmt"

	"encoding/base64"

	"log"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const rpcErrorKey = "rpc-error"

func (e *Error) Error() string {
	return e.Message
}

func Errorf(code codes.Code, userErrorCodes Error_UserErrorCode, temporary bool, msg string, args ...interface{}) error {
	return &Error{
		Code:          int64(code),
		Message:       fmt.Sprintf(msg, args...),
		Temporary:     temporary,
		UserErrorCode: userErrorCodes,
	}
}

func MashalError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	rerr, ok := err.(*Error)
	if !ok {
		return err
	}

	if pb, err := proto.Marshal(rerr); err == nil {
		md := metadata.Pairs(rpcErrorKey, base64.StdEncoding.EncodeToString(pb))
		if err := grpc.SetTrailer(ctx, md); err != nil {
			log.Println(err)
		}
	}

	return status.Errorf(codes.Code(rerr.Code), rerr.Message)
}

func UnmarshalError(err error, md metadata.MD) *Error {
	vals, ok := md[rpcErrorKey]
	if !ok {
		return nil
	}
	if len(vals) < 1 {
		return nil
	}

	buf, err := base64.StdEncoding.DecodeString(vals[0])
	if err != nil {
		return nil
	}

	rerr := &Error{}
	if err := proto.Unmarshal(buf, rerr); err != nil {
		return nil
	}

	return rerr
}
