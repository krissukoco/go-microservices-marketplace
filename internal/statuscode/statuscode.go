package statuscode

//  Codes consists of 5-digit non-negative integers

const (
	OK int = 0
	// 10000 - 19999: User & Auth Related
	Unauthorized           int = 10001
	Forbidden              int = 10002
	TokenGeneral           int = 10010
	TokenMissing           int = 10011
	TokenMalformed         int = 10012
	TokenExpired           int = 10013
	TokenInvalid           int = 10014
	EmailOrPasswordInvalid int = 10021
	TooManyRequests        int = 19001
	// 30000 - 39999: Resource Related
	ResourceGeneral   int = 30000
	ResourceForbidden int = 30003
	ResourceNotFound  int = 30004
	// 40000 - 49999: Request Validations
	UnparsableBody int = 40000
	// 90000 - 99999: Server Errors
	ServerUnknown      int = 90000
	ServiceUnavailable int = 90003
)

var HTTPCodeMap = map[int][]int{
	200: {OK},
	400: {EmailOrPasswordInvalid},
	401: {Unauthorized, TokenGeneral, TokenMissing, TokenMalformed, TokenExpired, TokenInvalid},
	403: {Forbidden, ResourceForbidden},
	404: {ResourceNotFound},
	422: {UnparsableBody},
	429: {TooManyRequests},
	500: {ServerUnknown},
	503: {ServiceUnavailable},
}

func HTTP(code int) int {
	for k, v := range HTTPCodeMap {
		for _, c := range v {
			if c == code {
				return k
			}
		}
	}
	return 400
}

func Message(code int) string {
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
	case ServerUnknown:
		return "Internal server error"
	case ServiceUnavailable:
		return "Service is unavailable"
	default:
		return "Unknown error"
	}
}
