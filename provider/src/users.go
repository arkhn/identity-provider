package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (env *Env) handleSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &User{}
	json.NewDecoder(r.Body).Decode(user)

	createdUser, err := env.db.AddUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(createdUser)
}
