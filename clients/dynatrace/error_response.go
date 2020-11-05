package dynatrace

import (
	"fmt"
	"strings"
)

type errorResponse struct {
	Details errorDetails `json:"error"`
}

func (e errorResponse) Error() string {
	return e.Details.Error()
}

type errorDetails struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	StatusMessage string `json:"-"`
}

func (e errorDetails) Error() string {
	msg := strings.TrimSpace(e.Message)
	if len(msg) == 0 {
		msg = strings.TrimSpace(e.StatusMessage)
	}
	if len(msg) == 0 {
		msg = "no error message available"
	}
	return fmt.Sprintf("[%d] %s", e.Code, msg)
}
