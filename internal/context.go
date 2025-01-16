package main

import (
	"context"
	"net/http"

	"github.com/luthfiaed/owapp-be/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		/*
			this function should only be called when we logically expect user object to exist in the context
			hence we panic when there's none
		*/
		panic("missing user value in request context")
	}
	return user
}
