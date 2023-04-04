package ecommerce

import (
	"encoding/json"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/user"
	"io"
	"net/http"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		logger.Errorf("can't process request body", err)
		w.WriteHeader(http.StatusBadRequest)
		//fb.errFn("failed loading body: %+s", err)
		//return it, http.StatusInternalServerError, errors.NewNotValid(err, "unable to read request body")
		return
	}
	var dto user.UserRequest
	err = json.Unmarshal(body, &dto)

	if err != nil {
		logger.Errorf("can't process request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = userService.NewUser(dto)
	if err != nil {
		logger.Errorf("can't add new user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
func addUserHandler() processing.ActivityHandlerFn {
	return func(receivedIn vocab.IRI, r *http.Request) (vocab.Item, int, error) {
		var it vocab.Item
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil || len(body) == 0 {
			logger.Errorf("can't process request body", err)
			return it, http.StatusBadRequest, err
		}
		var dto user.UserRequest
		err = json.Unmarshal(body, &dto)

		if err != nil {
			logger.Errorf("can't process request body", err)
			return it, http.StatusBadRequest, err
		}

		u, err := userService.NewUser(dto)
		if err != nil {
			logger.Errorf("can't add new user", err)
			return it, http.StatusBadRequest, err
		}
		return u, http.StatusCreated, nil
	}
}
*/
