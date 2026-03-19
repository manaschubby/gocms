package httpTransport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type output struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func Ok(e echo.Context, v any) error {
	bytes, ok := v.([]byte)
	if ok {
		return e.JSONBlob(http.StatusOK, bytes)
	}
	return e.JSON(http.StatusOK, output{
		Status: http.StatusOK,
		Data:   v,
	})
}

func Err(e echo.Context, code int, v any) error {
	return e.JSON(code, output{
		Status: code,
		Error:  http.StatusText(code),
		Data:   v,
	})
}

func ErrWithMsg(e echo.Context, code int, msg string, v any) error {
	return e.JSON(code, output{
		Status: code,
		Error:  msg,
		Data:   v,
	})
}
