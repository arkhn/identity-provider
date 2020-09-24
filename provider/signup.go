package provider

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"main/users"
)

func (ctx *Provider) GetSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderTemplate(w, "signup.html", nil)
}

func (ctx *Provider) PostSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse form").Error(), http.StatusBadRequest)
		return
	}

	user := &users.User{
		Name:     r.Form.Get("name"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	createdUser, err := ctx.Db.AddUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplate(w, "signup_success.html", createdUser)
}
