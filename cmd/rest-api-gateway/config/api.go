package config

type ApiConfigSchema struct {
	UserServiceUrl    string
	ChatServiceUrl    string
	ProductServiceUrl string
}

var Api *ApiConfigSchema

func initApiConfig() {
	Api = &ApiConfigSchema{
		// TODO: Put URLs in env vars
		UserServiceUrl:    "localhost:11000",
		ChatServiceUrl:    "localhost:12000",
		ProductServiceUrl: "localhost:13000",
	}

}
