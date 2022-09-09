package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerguard/src/oauth"
	"github.com/shivanshkc/ledgerguard/src/utils/errutils"
	"github.com/shivanshkc/ledgerguard/src/utils/httputils"

	"github.com/gorilla/mux"
)

// Providers
var (
	providerGoogle = oauth.NewProviderGoogle()
	// More provider implementations can be added here.
)

// ProviderMap maps the provider ID to the actual provider object.
var ProviderMap = map[string]oauth.Provider{
	providerGoogle.ID(): providerGoogle,
	// More provider implementations can be added here.
}

// AuthHandler is the HTTP handler for the authentication API.
func AuthHandler(writer http.ResponseWriter, req *http.Request) {
	// Prerequisites.
	ctx := req.Context()

	// Extracting the providerID path parameter. This is the identification os the SSO provider. Example: Google.
	providerID := mux.Vars(req)["provider_id"]
	// Extracting the stateRedirectURI. This will be used as the state parameter in the provider's auth API.
	state := req.URL.Query().Get("redirect_uri")

	// Validating the state param.
	if err := checkState(state); err != nil {
		errHTTP := errutils.BadRequest().WithReasonError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Checking if the specified provider exists. This also acts as a validation for the provider ID.
	provider, exists := ProviderMap[providerID]
	if !exists {
		errHTTP := errutils.ProviderNotFound()
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Obtaining the provider's redirect URI.
	providerRedirectURI, err := provider.GetRedirectURI(ctx, state)
	if err != nil {
		errHTTP := errutils.ToHTTPError(err)
		httputils.Write(writer, errHTTP.Status, nil, errHTTP)
		return // nolint:wsl // Allowing return statement cuddling.
	}

	// Response headers to actually make the redirection work.
	headers := map[string]string{"Location": providerRedirectURI}
	// Writing the final response.
	httputils.Write(writer, http.StatusFound, headers, nil)
}
