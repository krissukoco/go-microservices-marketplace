package config

import "os"

type ApiConfigSchema struct {
	ClientApiKey      string
	UserServiceUrl    string
	ChatServiceUrl    string
	ListingServiceUrl string
}

var Api *ApiConfigSchema

func initApiConfig() {
	clientApiKey, ok := os.LookupEnv("CLIENT_API_KEY")
	if !ok {
		panic("CLIENT_API_KEY not set")
	}
	Api = &ApiConfigSchema{
		ClientApiKey: clientApiKey,
		// TODO: Put URLs in env vars
		UserServiceUrl:    "localhost:11000",
		ChatServiceUrl:    "localhost:12000",
		ListingServiceUrl: "localhost:13000",
	}

}
