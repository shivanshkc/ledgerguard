package oauth

import (
	"context"

	"github.com/shivanshkc/ledgerguard/src/models"
)

// Provider represents an OAuth provider. Example: Google, Facebook.
type Provider interface {
	// ID provides the identifier of the provider.
	ID() string

	// GetRedirectURI provides the URI to which the client shall be redirected for authentication with the provider.
	GetRedirectURI(ctx context.Context, state string) (string, error)
	// GetIdentityToken fetches the identity token using the provided code.
	GetIdentityToken(ctx context.Context, code string) (string, error)
	// GetUserInfo provides the user's info from the identity token.
	GetUserInfo(ctx context.Context, identityToken string) (*models.OAuthUser, error)
}
