package middleware

import (
	"context"
	auth2 "github.com/go-ap/auth"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"net/http"
)

var AuthService *auth2.Server

// CheckToken -  reads the Authorization header and load an actor. The actor saved in request context
// If actor not found then http.StatusForbidden returned to the client.
func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
		ctx := r.Context()
		actor, err := AuthService.LoadActorFromAuthHeader(r)
		if err != nil {
			errors.HandleError(err).ServeHTTP(w, r)
			return
		}
		if actor.ID == auth2.AnonymousActor.GetID() {
			err = errors.Unauthorizedf("request unauthorised")
			errors.HandleError(err).ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, common.AuthActorKey{}, actor)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
