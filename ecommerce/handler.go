package ecommerce

import (
	"encoding/json"
	"github.com/go-ap/fedbox/ecommerce/user"
	"io"
	"net/http"
)

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		//fb.errFn("failed loading body: %+s", err)
		//return it, http.StatusInternalServerError, errors.NewNotValid(err, "unable to read request body")
	}
	var dto user.UserRequest
	err = json.Unmarshal(body, &dto)

	if err != nil {

	}

	_, err = userService.NewUser(dto)
	if err != nil {

	}
}
