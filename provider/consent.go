package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// After pressing "click here", the Authorize Code flow is performed and the user is redirected to Hydra. Next, Hydra
// validates the consent request (it's not valid yet) and redirects us to the consent endpoint which we set with `CONSENT_URL=http://localhost:4445/consent`.
func (ctx *AuthContext) GetConsent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the consent requerst id from the query.
	challenge := r.URL.Query().Get("consent_challenge")

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
	}{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)

	if err != nil {
		http.Error(w, errors.Wrap(err, "Could not parse consent request info").Error(), http.StatusInternalServerError)
	}

	requestedScopes := jsonResp.RequestedScopes
	// TODO do stuff with response

	// This helper checks if the user is already authenticated. If not, we
	// redirect them to the login endpoint.
	// user := authenticated(r)
	// if user == "" {
	// 	http.Redirect(w, r, "/login?consent="+consentRequestID, http.StatusFound)
	// 	return
	// }

	fillTemplate := struct {
		ConsentRequestID string
		ClientID         string
		RequestedScopes  []string
	}{
		ConsentRequestID: challenge,
		ClientID:         "app id", // TODO
		RequestedScopes:  requestedScopes,
	}

	renderTemplate(w, "consent.html", fillTemplate)
}

func (ctx *AuthContext) PostConsent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the consent requerst id from the query.
	challenge := r.URL.Query().Get("consent_challenge")

	params := url.Values{}
	params.Add("consent_challenge", challenge)

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

	putUrl := fmt.Sprintf("%s/accept?%s", ctx.HConf.ConsentRequestRoute, params.Encode())

	// TODO use session to add info about the current user
	session := SessionInfo{
		IdToken: IdTokenClaims{
			Name:  "bob",
			Email: "bob@arkhn.com",
		},
	}

	body := &BodyAcceptOAuth2Consent{
		GrantScope:               grantedScopes,
		GrantAccessTokenAudience: []string{"http://localhost:3002"}, // TODO
		Remember:                 false,
		RememberFor:              3600,
		Session:                  session,
	}

	putAndRedirect(putUrl, body, w, r, http.DefaultClient)
}
