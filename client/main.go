package main

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/upamune/grpc-error-sample/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func l(md metadata.MD, err error) {
	if err == nil {
		return
	}

	rerr := api.UnmarshalError(err, md)
	fmt.Printf(`
=== Error ===
Code:          %s
Message:       %s
Temporary:     %v
UserErrorCode: %s
==============
	`, codes.Code(rerr.Code).String(), rerr.Message, rerr.Temporary, rerr.UserErrorCode.String())
}

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewUserServiceClient(conn)

	req := &api.CreateUserRequest{
		Id:   "serizawa",
		Name: "Yu SERIZAWA",
	}
	_, err = client.Create(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Get(context.TODO(), &api.GetUserRequest{
		UserId: req.Id,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %v\n", res.User)

	var trailer metadata.MD

	_, err = client.Get(context.TODO(), &api.GetUserRequest{
		UserId: "yamakita",
	},
		grpc.Trailer(&trailer),
	)
	l(trailer, err)

	_, err = client.Create(context.TODO(), req,
		grpc.Trailer(&trailer),
	)
	l(trailer, err)

	_, err = client.Delete(context.TODO(), &api.DeleteUserRequest{
		UserId: req.Id,
	},
		grpc.Trailer(&trailer),
	)
	l(trailer, err)
}
