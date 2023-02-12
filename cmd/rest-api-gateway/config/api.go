package config

import "os"

type ApiConfigSchema struct {
	ClientApiKey string
}

var Api *ApiConfigSchema

func initApiConfig() {
	clientApiKey, ok := os.LookupEnv("CLIENT_API_KEY")
	if !ok {
		panic("CLIENT_API_KEY not set")
	}
	Api = &ApiConfigSchema{
		ClientApiKey: clientApiKey,
	}

}
