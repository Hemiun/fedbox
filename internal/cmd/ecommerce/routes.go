package ecommerce

import (
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ECommerceRoutes() func(router chi.Router) {
	return func(r chi.Router) {
		r.Use(middleware.CheckToken)
		r.Route("/ecommerce", func(r chi.Router) {
			r.Group(userRoutes())
			r.Group(productRoutes())
		})
	}
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: remove it
	w.Write([]byte(r.Method))
}

func userRoutes() func(router chi.Router) {
	return func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/", emptyHandler)
			r.Post("/", addUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Put("/", emptyHandler)
				r.Delete("/", deleteUserHandler)
				r.Get("/", emptyHandler)
			})

		})
	}
}

func productRoutes() func(router chi.Router) {
	return func(r chi.Router) {
		r.Route("/product", func(r chi.Router) {
			r.Get("/", emptyHandler)
			r.Get("/{id}", emptyHandler)
			r.Post("/", emptyHandler)
			r.Put("/", emptyHandler)
		})
	}
}
