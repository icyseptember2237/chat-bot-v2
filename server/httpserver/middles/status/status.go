package status

import "net/http"

type Status struct {
	Code    int
	Message string
}

func (s *Status) Error() string {
	return s.Message
}

func (s *Status) GetCode() int {
	return s.Code
}

// Error returns an error representing c and msg.  If c is OK, returns nil.
func Error(c int, msg string) error {
	return &Status{
		Code:    c,
		Message: msg,
	}
}

func GetCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err := err.(type) {
	case *Status:
		return err.Code
	default:
		return http.StatusInternalServerError
	}
}
