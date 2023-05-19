package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authPb "github.com/krissukoco/go-microservices-marketplace/proto/auth"
)

type AuthClient struct {
	Url         string
	DialOptions []grpc.DialOption
}

var Auth *AuthClient

func loadAuthTlsCreds() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("../../certificates/auth-ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add auth CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func SetupAuthClient() error {
	dialOptions := []grpc.DialOption{}

	// TLS credentials
	tlsCredentials, err := loadAuthTlsCreds()
	if err != nil {
		return err
	}
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(tlsCredentials))

	Auth = &AuthClient{Url: "localhost:11000"}
	Auth.DialOptions = dialOptions

	log.Println("Dial options: ", Auth)

	return nil
}

// =============================================== RPC ===============================================

func (c *AuthClient) Register(req *authPb.RegisterRequest) (*authPb.RegisterResponse, error) {
	conn, err := grpc.Dial(c.Url, c.DialOptions...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := authPb.NewAuthServiceClient(conn)
	res, err := client.RegisterUser(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AuthClient) Login(req *authPb.LoginRequest) (*authPb.LoginResponse, error) {
	conn, err := grpc.Dial(c.Url, c.DialOptions...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := authPb.NewAuthServiceClient(conn)
	res, err := client.Login(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AuthClient) Refresh(req *authPb.RefreshRequest) (*authPb.RefreshResponse, error) {
	conn, err := grpc.Dial(c.Url, c.DialOptions...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := authPb.NewAuthServiceClient(conn)
	res, err := client.Refresh(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
