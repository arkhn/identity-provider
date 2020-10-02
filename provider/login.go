package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/common/env"
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
	_, err = http.Get(getUrl)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while fetching login request info from hydra").Error(), http.StatusInternalServerError)
	}

	// TODO do stuff with response
	// var jsonResp interface{}
	// err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	// if err != nil {
	// 	http.Error(w, errors.Wrap(err, "Could not parse login request info").Error(), http.StatusInternalServerError)
	// }

	// log.Println(jsonResp)

	fillTemplate := struct {
		ConsentChallenge string
		RootURL          string
	}{
		ConsentChallenge: challenge,
		RootURL:          env.Getenv("ROOT_URL", ""),
	}
	renderTemplate(w, "login.html", fillTemplate)
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

	// Let's create a session where we store the user id. We can ignore errors from the session store
	// as it will always return a session!
	session, _ := ctx.Store.Get(r, sessionName)
	session.Values["userName"] = user.Name
	session.Values["userEmail"] = user.Email
	// Store the session in the cookies
	if err := ctx.Store.Save(r, w, session); err != nil {
		http.Error(w, errors.Wrap(err, "Could not persist cookies").Error(), http.StatusBadRequest)
		return
	}

	putUrl := fmt.Sprintf("%s/accept?%s", ctx.HConf.LoginRequestRoute, params.Encode())

	// TODO properly fill body
	body := &BodyAcceptOAuth2Login{
		// Acr:         "..",
		Remember:    false,
		RememberFor: 3600,
		Subject:     user.Email,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", putUrl, bytes.NewBuffer(jsonBody))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while accepting login request").Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while accepting login request").Error(), http.StatusInternalServerError)
	}

	jsonResp := RedirectResp{}
	err = json.Unmarshal(b, &jsonResp)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while accepting login request").Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, jsonResp.RedirectTo, http.StatusFound)
}
