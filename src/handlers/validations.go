package handlers

import (
	"errors"
	"net/url"
)

// ErrInvalidRedirectURI is the error returned when the state redirect-uri is invalid.
var ErrInvalidRedirectURI = errors.New("invalid redirect uri")

// checkState checks if the specified redirect URI is valid.
func checkState(uri string) error {
	// Parsing the URI to check if it's valid.
	if _, err := url.ParseRequestURI(uri); err != nil {
		return ErrInvalidRedirectURI
	}

	return nil
}
