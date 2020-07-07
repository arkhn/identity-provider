package provider

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"main/users"
)

func (ctx *AuthContext) HandleSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	createdUser, err := ctx.Db.AddUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
