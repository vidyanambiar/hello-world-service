// Copyright Red Hat

package routes

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/identitatem/idp-configs-api/pkg/services"
)

// Routing for operations on Auth Realms
func MakeRouterForAuthRealms(sub chi.Router) {
	sub.Get("/", services.GetAuthRealmsForAccount)
	sub.Post("/", services.CreateAuthRealmForAccount)
	sub.Route("/{id}", func(r chi.Router) {
		r.Use(services.AuthRealmCtx)
		r.Get("/", services.GetAuthRealmByID)
		r.Put("/", services.UpdateAuthRealmByID)
		r.Delete("/", services.DeleteAuthRealmByID)
	})
}