package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krissukoco/go-microservices-marketplace/internal/statuscode"
)

type APIResponseError struct {
	OK      bool   `json:"ok"`
	Code    int64  `json:"code"`
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

func NewErrorResponse(code int64, msg string) *APIResponseError {
	return &APIResponseError{
		OK:      false,
		Code:    code,
		Message: msg,
	}
}

func APIOkWithData(c *fiber.Ctx, data interface{}) error {
	return c.Status(200).JSON(&APIResponseOK{
		OK:   true,
		Data: data,
	})
}

func APIErrorFromCode(c *fiber.Ctx, code int64, message ...string) error {
	// TODO: include params
	msg := statuscode.Message(code)
	// Override message
	if len(message) > 0 {
		msg = message[0]
	}
	httpCode := statuscode.HTTP(code)
	return c.Status(httpCode).JSON(&APIResponseError{
		OK:      false,
		Code:    code,
		Message: msg,
	})
}
