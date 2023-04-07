package ecommerce

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/user"
	json "github.com/go-ap/jsonld"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const MIMEApplicationJSON = "application/json"

type UserService interface {
	CreateUser(caller vocab.Actor, ur user.UserDTO) (vocab.Item, error)
	FindUser(caller vocab.Actor, userID string) (*user.UserDTO, error)
	DeleteUser(caller vocab.Actor, userID string) error
	UpdateUser(caller vocab.Actor, userID string, ur user.UserDTO) (vocab.Item, error)
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		err = errors.NewBadRequest(err, "not empty body expected")
		logger.Errorf("not empty body expected", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	defer r.Body.Close()

	var dto user.UserDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		err = errors.NewBadRequest(err, "can't process request body")
		logger.Errorf("can't process request body", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	it, err := userService.CreateUser(callerActor, dto)
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

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		err = errors.NewBadRequest(err, "userID not passed")
		logger.Errorf("userID not passed", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	usr, err := userService.FindUser(callerActor, userID)
	if err != nil {
		logger.Errorf("can't find user", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	data, err := json.Marshal(usr)
	if err != nil {
		logger.Errorf("can't marshall response", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Type", MIMEApplicationJSON)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		err = errors.NewBadRequest(err, "userID not passed")
		logger.Errorf("userID not passed", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	err = userService.DeleteUser(callerActor, userID)
	if err != nil {
		logger.Errorf("can't delete user", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	callerActor, ok := r.Context().Value(common.AuthActorKey{}).(vocab.Actor)
	if !ok {
		err = errors.NewBadRequest(err, "can't get actor from context")
		logger.Errorf("can't get actor from context", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		err = errors.NewBadRequest(err, "userID not passed")
		logger.Errorf("userID not passed", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		err = errors.NewBadRequest(err, "not empty body expected")
		logger.Errorf("not empty body expected", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}
	defer r.Body.Close()

	var dto user.UserDTO
	err = json.Unmarshal(body, &dto)
	if err != nil {
		err = errors.NewBadRequest(err, "can't process request body")
		logger.Errorf("can't process request body", err)
		errors.HandleError(err).ServeHTTP(w, r)
		return
	}

	it, err := userService.UpdateUser(callerActor, userID, dto)
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
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("can't write response", err)
	}
}
