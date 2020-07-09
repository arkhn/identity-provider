package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func (ctx *Provider) GetLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	challenge, err := parseChallengeFromRequest(r, "login_challenge")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Get info about current flow
	params := url.Values{}
	params.Add("login_challenge", challenge)

	getUrl := fmt.Sprintf("%s?%s", ctx.HConf.LoginRequestRoute, params.Encode())
	resp, err := http.Get(getUrl)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while fetching login request info from hydra").Error(), http.StatusInternalServerError)
	}

	var jsonResp interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse login request info").Error(), http.StatusInternalServerError)
	}

	log.Println(jsonResp)
	// TODO do stuff with response

	renderTemplate(w, "login.html", challenge)
}

func (ctx *Provider) PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	challenge, err := parseChallengeFromRequest(r, "login_challenge")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	params := url.Values{}
	params.Add("login_challenge", challenge)

	// Parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse form").Error(), http.StatusBadRequest)
		return
	}

	// Check the user's credentials
	user, err := ctx.Db.FindUser(r.Form.Get("email"), r.Form.Get("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(user)

	// // Let's create a session where we store the user id. We can ignore errors from the session store
	// // as it will always return a session!
	// session, _ := store.Get(r, sessionName)
	// session.Values["user"] = "buzz-lightyear"

	// // Store the session in the cookie
	// if err := store.Save(r, w, session); err != nil {
	// 	http.Error(w, errors.Wrap(err, "Could not persist cookie").Error(), http.StatusBadRequest)
	// 	return
	// }

	// Redirect the user back to the consent endpoint. In a normal app, you would probably
	// add some logic here that is triggered when the user actually performs authentication and is not
	// part of the consent flow.
	putUrl := fmt.Sprintf("%s/accept?%s", ctx.HConf.LoginRequestRoute, params.Encode())

	// TODO properly fill body
	body := &BodyAcceptOAuth2Login{
		Acr:         "..",
		Remember:    false,
		RememberFor: 3600,
		Subject:     "bob",
	}

	putAndRedirect(putUrl, body, w, r, http.DefaultClient)
}
