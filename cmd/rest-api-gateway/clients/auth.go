package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthClient struct {
	Url         string
	DialOptions []grpc.DialOption
}

var Auth *AuthClient

func loadAuthTlsCreds() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("../../../certificates/auth-ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add auth CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
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
	Auth.DialOptions = dialOptions

	return nil
}
