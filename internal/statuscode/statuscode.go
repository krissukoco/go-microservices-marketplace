package statuscode

import (
	"fmt"
	"strconv"
	"strings"
)

//  Codes consists of 5-digit non-negative integers

const (
	OK int64 = 0
	// 10000 - 19999: User & Auth Related
	Unauthorized           int64 = 10001
	Forbidden              int64 = 10002
	TokenGeneral           int64 = 10010
	TokenMissing           int64 = 10011
	TokenMalformed         int64 = 10012
	TokenExpired           int64 = 10013
	TokenInvalid           int64 = 10014
	EmailOrPasswordInvalid int64 = 10021
	EmailAlreadyRegistered int64 = 10022
	TooManyRequests        int64 = 19001
	// 30000 - 39999: Resource Related
	ResourceGeneral   int64 = 30000
	ResourceForbidden int64 = 30003
	ResourceNotFound  int64 = 30004
	// --> 31000 - 31999: Product related
	ProductNotFound           int64 = 31004
	ProductTitleTooShort      int64 = 31111
	ProductTitleTooLong       int64 = 31112
	ProductDescriptionTooLong int64 = 31121
	ProductCategoryEmpty      int64 = 31131
	ProductCategoryInvalid    int64 = 31132
	ProductPriceInvalid       int64 = 31141
	ProductPriceTooLow        int64 = 31142
	ProductTagsTooLong        int64 = 31151
	// 40000 - 49999: Request Validations
	BadRequest        int64 = 40000
	PasswordMalformed int64 = 40001
	UnparsableBody    int64 = 42000
	// 90000 - 99999: Server & Unknown Errors
	ServerError        int64 = 90000
	ServiceUnavailable int64 = 90003
	UnknownError       int64 = 99999
)

var HTTPCodeMap = map[int][]int64{
	200: {OK},
	400: {BadRequest, EmailOrPasswordInvalid, EmailAlreadyRegistered, UnknownError},
	401: {Unauthorized, TokenGeneral, TokenMissing, TokenMalformed, TokenExpired, TokenInvalid},
	403: {Forbidden, ResourceForbidden},
	404: {ResourceNotFound},
	422: {UnparsableBody},
	429: {TooManyRequests},
	500: {ServerError},
	503: {ServiceUnavailable},
}

func HTTP(code int64) int {
	for k, v := range HTTPCodeMap {
		for _, c := range v {
			if c == code {
				return k
			}
		}
	}
	return 400
}

func Message(code int64) string {
	switch code {
	case OK:
		return "OK"
	case Unauthorized:
		return "Unauthorized"
	case Forbidden:
		return "Forbidden"
	case TokenGeneral:
		return "Token error"
	case TokenMissing:
		return "Token is missing"
	case TokenMalformed:
		return "Token is malformed"
	case TokenExpired:
		return "Token is expired"
	case TokenInvalid:
		return "Token is invalid"
	case EmailOrPasswordInvalid:
		return "Email or password is invalid"
	case TooManyRequests:
		return "Too many requests"
	case ResourceGeneral:
		return "Resource error"
	case ResourceForbidden:
		return "Resource is forbidden"
	case ResourceNotFound:
		return "Resource is not found"
	case UnparsableBody:
		return "Unparsable body"
	case ServerError:
		return "Internal server error"
	case ServiceUnavailable:
		return "Service is unavailable"
	default:
		return "Unknown error"
	}
}

func ParseGrpcErrMsg(msg string) (int64, string) {
	fmt.Println("GRPC Message: ", msg)
	code := UnknownError
	message := "Unknown error"
	split := strings.Split(msg, "__")
	if len(split) != 2 {
		return code, message
	}
	intCode, err := strconv.Atoi(split[0])
	if err != nil {
		return code, message
	}
	code = int64(intCode)
	message = split[1]
	return code, message
}

func StandardErrorMessage(code int64, message string) string {
	return fmt.Sprintf("%d__%s", code, message)
}
