package main

import (
	"expvar"
	"net/http"
)

func (app *application) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	// Use the requirePermission() middleware on each of the /v1/movies** endpoints,
	// passing in the required permission code as the first parameter
	mux.HandleFunc("/v1/home", app.requirePermisson("movies:read", http.HandlerFunc(app.listMoviesHandler)))
	mux.HandleFunc("/v1/movies", app.requirePermisson("movies:write", http.HandlerFunc(app.createMovieHandler)))
	mux.HandleFunc("/v1/onemovies", app.requirePermisson("movies:read", http.HandlerFunc(app.showMovieHandler)))
	mux.HandleFunc("/v1/updatemovies", app.requirePermisson("movies:write", http.HandlerFunc(app.updateMovieHandler)))
	mux.HandleFunc("/v1/delete", app.requirePermisson("movies:write", http.HandlerFunc(app.deleteMovieHandler)))
	// Add the route for the POST /v1/users endpoint.
	mux.HandleFunc("/v1/users", app.registerUserHandler)
	mux.HandleFunc("/v1/users/activated", app.activateUserHandler)
	mux.HandleFunc("/v1/users/password", app.updateUserPasswordHandler)
	mux.HandleFunc("/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	mux.HandleFunc("/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	// Reagister a new Get /debug/vars endpont pointing to the expvar handler
	mux.Handle("/v1/metrics", expvar.Handler())
	return app.metrics(app.recoverPanic(app.enableCORS(app.authenticate(mux))))
}
