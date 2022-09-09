package oauth

import (
	"context"
	"fmt"

	"github.com/shivanshkc/ledgerguard/src/models"
)

// providerGoogle implements the Provider interface using Google's sign-in.
type providerGoogle struct{}

func (p *providerGoogle) ID() string {
	return "google"
}

func (p *providerGoogle) GetRedirectURI(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(
		"%s?scope=%s&include_granted_scopes=true&response_type=code&redirect_uri=%s&client_id=%s",
		"", // TODO: Google's auth endpoint.
		"", // TODO: Google's scopes.
		fmt.Sprintf("%s/api/auth/%s/callback", "" /* TODO: Own base URL */, p.ID()),
		"", // TODO: Google's client ID.
	), nil
}

func (p *providerGoogle) GetIdentityToken(ctx context.Context, code string) (string, error) {
	panic("implement me")
}

func (p *providerGoogle) GetUserInfo(ctx context.Context, identityToken string) (*models.OAuthUser, error) {
	panic("implement me")
}
