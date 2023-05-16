package config

type ApiConfigSchema struct {
	UserServiceUrl    string
	ChatServiceUrl    string
	ListingServiceUrl string
}

var Api *ApiConfigSchema

func initApiConfig() {
	Api = &ApiConfigSchema{
		// TODO: Put URLs in env vars
		UserServiceUrl:    "localhost:11000",
		ChatServiceUrl:    "localhost:12000",
		ListingServiceUrl: "localhost:13000",
	}

}
