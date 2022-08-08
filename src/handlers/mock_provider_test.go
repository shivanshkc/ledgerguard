package handlers_test

import (
	"context"

	"github.com/shivanshkc/ledgerguard/src/models"
)

// mockProvider is a mock implementation of an OAuth provider.
type mockProvider struct {
	// providerID can be used to mock the value returned by the ID method.
	providerID string

	// errRedirectURI can be used to mock the error returned by the GetRedirectURI method.
	errRedirectURI error
	// redirectURI can be used to mock the uri returned by the GetRedirectURI method.
	redirectURI string
}

func (m *mockProvider) ID() string {
	return m.providerID
}

func (m *mockProvider) GetRedirectURI(ctx context.Context, state string) (string, error) {
	// If there's a mock error present, it is returned.
	if m.errRedirectURI != nil {
		return "", m.errRedirectURI
	}

	return m.redirectURI, nil
}

func (m *mockProvider) GetIdentityToken(ctx context.Context, code string) (string, error) {
	panic("implement me")
}

func (m *mockProvider) GetUserInfo(ctx context.Context, identityToken string) (*models.OAuthUser, error) {
	panic("implement me")
}
