package errutils

import (
	"net/http"
)

// ProviderNotFound is for requests that require an OAuth provider that does not exist.
func ProviderNotFound() *HTTPError {
	return &HTTPError{Status: http.StatusNotFound, Code: "PROVIDER_NOT_FOUND"}
}
