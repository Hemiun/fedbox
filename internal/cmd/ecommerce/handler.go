package ecommerce

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/user"
	json "github.com/go-ap/jsonld"
	"io"
	"net/http"
)

type UserService interface {
	AddUser(ur user.UserRequest, actor vocab.Actor) (vocab.Item, error)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		err = errors.NewBadRequest(err, "not empty body expected")
		logger.Errorf("not empty body expected", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	// TODO: check unmarshall (json -> jsonld)
	var dto user.UserRequest
	err = json.Unmarshal(body, &dto)
	if err != nil {
		err = errors.NewBadRequest(err, "can't process request body")
		logger.Errorf("can't process request body", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	actor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	it, err := userService.AddUser(dto, actor)
	if err != nil {
		logger.Errorf("can't add new user", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	var data []byte

	if data, err = vocab.MarshalJSON(it); err != nil {
		logger.Errorf("can't marshall response", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	w.Header().Set("Location", it.GetLink().String())
	w.Header().Set("Content-Type", json.ContentType)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}
