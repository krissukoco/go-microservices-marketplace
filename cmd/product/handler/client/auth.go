package client

import (
	"context"

	authPb "github.com/krissukoco/go-microservices-marketplace/proto/auth"
	"google.golang.org/grpc"
)

type AuthClient struct {
	Url string

	DialOptions []grpc.DialOption
}

func (c *AuthClient) GetUserIdByToken(ctx context.Context, token string) (string, error) {
	conn, err := grpc.Dial(c.Url, c.DialOptions...)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	client := authPb.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), &authPb.RefreshRequest{
		Token: token,
	})
	if err != nil {
		return "", err
	}
	return res.Id, nil
}
