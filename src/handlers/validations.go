package handlers

import (
	"fmt"
	"net/url"
)

// checkProviderID checks if the specified provider ID is valid.
func checkProviderID(providerID string) error {
	// We don't need any validations here.
	return nil
}

// checkRedirectURI checks if the specified redirect URI is valid.
func checkRedirectURI(uri string) error {
	// Parsing the URI to check if it's valid.
	if _, err := url.ParseRequestURI(uri); err != nil {
		return fmt.Errorf("invalid redirect uri: %w", err)
	}

	return nil
}
