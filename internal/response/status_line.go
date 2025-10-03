package response

import (
	"fmt"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	NotFound            StatusCode = 404
	InternalServerError StatusCode = 500
)

func statusToString(statusCode StatusCode) (string, error) {
	switch statusCode {
	case OK:
		return "OK", nil
	case BadRequest:
		return "Bad Request", nil
	case NotFound:
		return "Not Found", nil
	case InternalServerError:
		return "Internal Server Error", nil
	default:
		return "", fmt.Errorf("Unknown status code: %q", statusCode)
	}
}
