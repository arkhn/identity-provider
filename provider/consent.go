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

// After pressing "click here", the Authorize Code flow is performed and the user is redirected to Hydra. Next, Hydra
// validates the consent request (it's not valid yet) and redirects us to the consent endpoint which we set with `CONSENT_URL=http://localhost:4445/consent`.
func (ctx *Provider) GetConsent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the consent requerst id from the query.
	challenge, err := parseChallengeFromRequest(r, "consent_challenge")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Get info about current flow
	params := url.Values{}
	params.Add("consent_challenge", challenge)

	getUrl := fmt.Sprintf("%s?%s", ctx.HConf.ConsentRequestRoute, params.Encode())
	resp, err := http.Get(getUrl)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while fetching consent request info from hydra").Error(), http.StatusInternalServerError)
	}

	jsonResp := struct {
		RequestedScopes []string `json:"requested_scope"`
		Client          struct {
			ClientID   string `json:"client_id"`
			ClientName string `json:"client_name"`
		} `json:"client"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse consent request info").Error(), http.StatusInternalServerError)
	}

	requestedScopes := jsonResp.RequestedScopes
	clientID := jsonResp.Client.ClientID

	// NOTE we'll skip the consent phase of the flow for now because:
	// - we don't really use scopes for now
	// - we expect all the scopes to be accepted for each client
	// - we only use first party clients
	if true { // TODO find a way to determine which clients are first party
		ctx.grantScopes(requestedScopes, challenge, w, r)
	} else {
		// This helper checks if the user is already authenticated. If not, we
		// redirect them to the login endpoint.
		// user := authenticated(r)
		// if user == "" {
		// 	http.Redirect(w, r, "/login?consent="+consentRequestID, http.StatusFound)
		// 	return
		// }

		fillTemplate := struct {
			ConsentChallenge string
			ClientID         string
			RequestedScopes  []string
			RootURL          string
		}{
			ConsentChallenge: challenge,
			ClientID:         clientID,
			RequestedScopes:  requestedScopes,
			RootURL:          env.Getenv("ROOT_URL", ""),
		}

		renderTemplate(w, "consent.html", fillTemplate)
	}

}

func (ctx *Provider) PostConsent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	challenge, err := parseChallengeFromRequest(r, "consent_challenge")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Parse the HTTP form - required by Go.
	if err := r.ParseForm(); err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse form").Error(), http.StatusBadRequest)
		return
	}

	// Let's check which scopes the user granted.
	var grantedScopes = []string{}
	for key := range r.PostForm {
		// And add each scope to the list of granted scopes.
		grantedScopes = append(grantedScopes, key)
	}

	ctx.grantScopes(grantedScopes, challenge, w, r)
}

func (ctx *Provider) grantScopes(grantedScopes []string, consentChallenge string, w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Add("consent_challenge", consentChallenge)

	putUrl := fmt.Sprintf("%s/accept?%s", ctx.HConf.ConsentRequestRoute, params.Encode())

	// TODO use session to add info about the current user
	session := SessionInfo{
		IdToken: IdTokenClaims{
			Name:  "admin",
			Email: "admin@arkhn.com",
		},
	}

	body := &BodyAcceptOAuth2Consent{
		GrantScope:               grantedScopes,
		GrantAccessTokenAudience: []string{"http://localhost:3002"}, // TODO
		Remember:                 false,
		RememberFor:              3600,
		Session:                  session,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", putUrl, bytes.NewBuffer(jsonBody))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while accepting consent request").Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, errors.Wrap(err, "Error while accepting consent request").Error(), http.StatusInternalServerError)
	}

	jsonResp := RedirectResp{}
	json.Unmarshal(b, &jsonResp)

	http.Redirect(w, r, jsonResp.RedirectTo, http.StatusFound)
}
