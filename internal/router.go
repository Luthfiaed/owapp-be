package main

import (
	"net/http"
)

type router struct {
	mux         *http.ServeMux
	middlewares []middleware
}

/*
this `use` function serves as syntactic sugar
to more elegantly chain multiple middlewares
*/
func (r *router) use(mdw middleware) {
	r.middlewares = append(r.middlewares, mdw)
}

func (r *router) handle(pattern string, handler http.HandlerFunc) {
	h := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}
	r.mux.Handle(pattern, h)
}

func newRouter() *router {
	return &router{
		mux:         http.NewServeMux(),
		middlewares: []middleware{},
	}
}

func (app *application) InitHandler() *http.ServeMux {
	router := newRouter()

	router.use(app.authenticate)

	router.handle("GET /api/v1/healthcheck", app.healthcheckHandler)
	router.handle("GET /api/v1/users", app.requireLogin(app.getUserByUsername))
	router.handle("POST /api/v1/users", app.registerNewUser)
	router.handle("GET /api/v1/avatar/{filename}", app.getAvatar)
	router.handle("POST /api/v1/avatar", app.uploadAvatar)
	router.handle("GET /api/v1/products", app.requireLogin(app.getProducts))
	router.handle("GET /api/v1/products/{id}", app.requireLogin(app.getProductById))
	router.handle("POST /api/v1/products/review", app.requireLogin(app.createReview))
	router.handle("PATCH /api/v1/products/review", app.requireLogin(app.updateReview))
	router.handle("POST /api/v1/login", app.login)

	router.handle("/", app.notFoundResponse)

	return router.mux
}
