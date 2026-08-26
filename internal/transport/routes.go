package transport

import (
	"net/http"
	"repairdesk.local/internal/service"
)

func Routes(s *service.Service) http.Handler { return New(s) }
func MethodAllowed(method string) bool {
	return method == http.MethodGet || method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete
}
func ParseLimit(raw string) int {
	if raw == "" {
		return 25
	}
	n := 0
	for _, c := range raw {
		if c < '0' || c > '9' {
			return 25
		}
		n = n*10 + int(c-'0')
		if n > 100 {
			return 100
		}
	}
	if n == 0 {
		return 25
	}
	return n
}
func FriendlyStatus(code int) string {
	switch {
	case code >= 500:
		return "The repair desk is temporarily unavailable."
	case code == 404:
		return "We could not find that request."
	case code == 400:
		return "Please check the request details."
	}
	return "Request completed."
}
