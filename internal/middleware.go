package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/luthfiaed/owapp-be/internal/data"
	"github.com/pascaldekloe/jwt"
)

type middleware = func(http.HandlerFunc) http.HandlerFunc

func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			The "Vary" header indicates to any caches that the response may vary based on
			the value of the Authorization header in the request
		*/
		w.Header().Add("Vary", "Authorization")

		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		authzHeader := r.Header.Get("Authorization")
		if authzHeader == "" {
			r = app.contextSetUser(r, data.ANONYMOUS_USER)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authzHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.JwtSecret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != "owapp-be" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.AcceptAudience("owapp") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.users.GetByUsername(claims.Subject)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}

		/*
			VULNERABILITY POINT: there is no checking that the user who sent the request
			is the one that the token was made for
		*/
		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireLogin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
