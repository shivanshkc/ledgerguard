package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerguard/src/oauth"
	"github.com/shivanshkc/ledgerguard/src/utils/errutils"
	"github.com/shivanshkc/ledgerguard/src/utils/httputils"

	"github.com/gorilla/mux"
)

// providerMap maps the provider ID to the actual provider object.
var providerMap = map[string]oauth.Provider{}

// AuthHandler is the HTTP handler for the authentication API.
func AuthHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx := req.Context()

	// Extracting the providerID path parameter. This is the identification os the SSO provider. Example: Google.
	providerID := mux.Vars(req)["provider_id"]
	// Extracting the redirectURI. This is the URL that will be opened after authentication success/failure.
	redirectURI := req.URL.Query().Get("redirect_uri")

	// Validating the provider ID.
	if err := checkProviderID(providerID); err != nil {
		errHTTP := errutils.BadRequest().WithReasonError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Validating the redirect URI.
	if err := checkRedirectURI(redirectURI); err != nil {
		errHTTP := errutils.BadRequest().WithReasonError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Checking if the specified provider exists.
	provider, exists := providerMap[providerID]
	if !exists {
		errHTTP := errutils.ProviderNotFound()
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Obtaining the provider's redirect URI.
	providerRedirectURI, err := provider.GetRedirectURI(ctx, redirectURI)
	if err != nil {
		errHTTP := errutils.ToHTTPError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Response headers to actually make the redirection work.
	headers := map[string]string{"Location": providerRedirectURI}
	// Writing the final response.
	httputils.Write(writer, http.StatusOK, headers, nil)
}
