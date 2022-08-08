package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shivanshkc/ledgerguard/src/handlers"

	"github.com/shivanshkc/ledgerguard/src/configs"
	"github.com/shivanshkc/ledgerguard/src/logger"
	"github.com/shivanshkc/ledgerguard/src/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	// Prerequisites.
	ctx, conf, log := context.Background(), configs.Get(), logger.Get()

	// Logging the HTTP server details.
	log.Info(ctx, &logger.Entry{Payload: fmt.Sprintf("%s v%s http server starting at: %s",
		conf.Application.Name, conf.Application.Version, conf.HTTPServer.Addr)})

	// Starting the HTTP server.
	if err := http.ListenAndServe(conf.HTTPServer.Addr, handler()); err != nil {
		panic("failed to start http server:" + err.Error())
	}
}

// handler is responsible to handle all incoming HTTP traffic.
func handler() http.Handler {
	router := mux.NewRouter()

	// Attaching global middlewares.
	// Note that the order here is VERY important.
	router.Use(middlewares.Recovery)
	router.Use(middlewares.RequestContext)
	router.Use(middlewares.AccessLogger)
	router.Use(middlewares.CORS)

	// Client-facing authentication API.
	router.HandleFunc("/api/auth/{provider_id}", handlers.AuthHandler).
		Methods(http.MethodOptions, http.MethodGet)
	// SSO-provider-facing callback API.
	router.HandleFunc("/api/auth/{provider_id}/callback", handlers.AuthCallbackHandler).
		Methods(http.MethodOptions, http.MethodGet)

	// Get user API.
	router.HandleFunc("/api/user", handlers.GetUserHandler).
		Methods(http.MethodOptions, http.MethodGet)
	// Verify user API.
	router.HandleFunc("/api/user", handlers.VerifyUserHandler).
		Methods(http.MethodOptions, http.MethodHead)

	return router
}
