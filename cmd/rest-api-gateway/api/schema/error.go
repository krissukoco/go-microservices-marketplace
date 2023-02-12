package schema

// type ErrorCode int

// Error Codes consists of 5-digit non-negative integers
// 1st digit: 1 - Auth, 2 - Schema/Validation, 3 - Listing
const (
	ErrorOK int = 0
	// 10000 - 19999: User & Auth Related
	ErrorUnauthorized           int = 10001
	ErrorForbidden              int = 10002
	ErrorToken                  int = 10010
	ErrorTokenMissing           int = 10011
	ErrorTokenMalformed         int = 10012
	ErrorTokenExpired           int = 10013
	ErrorTokenInvalid           int = 10014
	ErrorEmailOrPasswordInvalid int = 10021
	// 30000 - 39999: Resource Related
	ErrorResourceGeneral   int = 30000
	ErrorResourceForbidden int = 30003
	ErrorResourceNotFound  int = 30004
	// 80000 - 89999: Request Validations
	ErrorUnparsableBody int = 80000
	// 90000 - 99999: Server Errors
	ErrorInternal           int = 90000
	ErrorServiceUnavailable int = 90003
)

type APIResponseError struct {
	OK      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type APIResponseOK struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data interface{}) *APIResponseOK {
	return &APIResponseOK{
		OK:   true,
		Data: data,
	}
}

func NewErrorResponse(code int, msg string) *APIResponseError {
	return &APIResponseError{
		OK:      false,
		Code:    code,
		Message: msg,
	}
}

func UnparsableBodyError() *APIResponseError {
	return NewErrorResponse(ErrorUnparsableBody, "Unparsable body")
}
