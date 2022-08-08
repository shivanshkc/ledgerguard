package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/shivanshkc/ledgerguard/src/handlers"
	"github.com/shivanshkc/ledgerguard/src/oauth"
	"github.com/shivanshkc/ledgerguard/src/utils/datautils"
	"github.com/shivanshkc/ledgerguard/src/utils/errutils"

	"github.com/gorilla/mux"
)

// nolint:funlen // Function is long due to table driven tests.
func TestAuthHandler(t *testing.T) {
	validProviderID := "google"
	stateRedirectURI := "https://my-site.com"
	providerRedirectURI := "https://provider-site.com"

	// Sugar for errors.
	errPNF := errutils.ProviderNotFound()
	errBad := errutils.BadRequest()
	errInt := errutils.InternalServerError()

	// Error that's returned by the mock provider in one of the cases.
	errProvider := errors.New("some error occurred")

	tests := []struct {
		// Request parameters.
		providerID       string
		stateRedirectURI string

		// Provider details.
		expectedProviderID string
		provider           oauth.Provider

		// Expected response.
		expectedStatus  int
		expectedHeaders http.Header
		expectedBody    map[string]interface{}
	}{
		// Success case.
		{
			providerID:         validProviderID,
			stateRedirectURI:   stateRedirectURI,
			expectedProviderID: validProviderID,
			provider:           &mockProvider{providerID: validProviderID, redirectURI: providerRedirectURI},
			expectedStatus:     http.StatusFound,
			expectedHeaders:    http.Header{"Location": []string{providerRedirectURI}},
			expectedBody:       nil,
		},
		// Unknown Provider ID case.
		{
			providerID:         validProviderID + "something",
			stateRedirectURI:   stateRedirectURI,
			expectedProviderID: validProviderID,
			provider:           &mockProvider{providerID: validProviderID, redirectURI: providerRedirectURI},
			expectedStatus:     http.StatusNotFound,
			expectedHeaders:    http.Header{},
			expectedBody:       map[string]interface{}{"code": errPNF.Code, "reason": errPNF.Reason},
		},
		// Invalid state-redirect-uri case.
		{
			providerID:         validProviderID,
			stateRedirectURI:   "",
			expectedProviderID: validProviderID,
			provider:           &mockProvider{providerID: validProviderID, redirectURI: providerRedirectURI},
			expectedStatus:     http.StatusBadRequest,
			expectedHeaders:    http.Header{},
			expectedBody:       map[string]interface{}{"code": errBad.Code, "reason": handlers.ErrState.Error()},
		},
		// Provider error case.
		{
			providerID:         validProviderID,
			stateRedirectURI:   stateRedirectURI,
			expectedProviderID: validProviderID,
			provider:           &mockProvider{providerID: validProviderID, errRedirectURI: errProvider},
			expectedStatus:     http.StatusInternalServerError,
			expectedHeaders:    http.Header{},
			expectedBody:       map[string]interface{}{"code": errInt.Code, "reason": errProvider.Error()},
		},
	}

	// Executing all tests.
	for _, test := range tests {
		// Putting the specified provider in the map.
		if test.expectedProviderID != "" && test.provider != nil {
			handlers.ProviderMap[test.expectedProviderID] = test.provider
		}

		// API route.
		target := fmt.Sprintf("/api/auth/provider_id?redirect_uri=%s", test.stateRedirectURI)
		// Dummy request.
		req := httptest.NewRequest(http.MethodGet, target, nil)
		// Setting the route params. Otherwise, the mux.Vars function won't work.
		req = mux.SetURLVars(req, map[string]string{"provider_id": test.providerID})

		// Dummy writer.
		writer := httptest.NewRecorder()

		// Executing the function to be tested.
		handlers.AuthHandler(writer, req)

		// Verifying response code.
		if writer.Code != test.expectedStatus {
			t.Errorf("expected response status to be: %d but got: %d", test.expectedStatus, writer.Code)
			return
		}
		// Verifying the Location header.
		if !reflect.DeepEqual(writer.Header().Get("Location"), test.expectedHeaders.Get("Location")) {
			t.Errorf("expected headers to be: %+v but got: %+v", test.expectedHeaders, writer.Header())
			return
		}

		// Converting the response body to map.
		var responseBody map[string]interface{}
		if err := datautils.AnyToAny(writer.Body, &responseBody); err != nil {
			t.Errorf("failed to verify response body: %+v", err)
			return
		}

		// Verifying the response body.
		if !reflect.DeepEqual(responseBody, test.expectedBody) {
			t.Errorf("expected body to be: %+v but got: %+v", test.expectedBody, responseBody)
			return
		}
	}
}
