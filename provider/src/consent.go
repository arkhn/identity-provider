package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// After pressing "click here", the Authorize Code flow is performed and the user is redirected to Hydra. Next, Hydra
// validates the consent request (it's not valid yet) and redirects us to the consent endpoint which we set with `CONSENT_URL=http://localhost:4445/consent`.
func handleConsent(hConf *hydraConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the consent requerst id from the query.
		challenge := r.URL.Query().Get("consent_challenge")

		// Get info about current flow
		params := url.Values{}
		params.Add("consent_challenge", challenge)

		getUrl := fmt.Sprintf("%s?%s", hConf.ConsentRequestRoute, params.Encode())
		resp, _ := http.Get(getUrl)

		b, _ := ioutil.ReadAll(resp.Body)
		var jsonResp interface{}
		json.Unmarshal(b, &jsonResp)
		log.Println(jsonResp)
		// TODO do stuff with response

		// This helper checks if the user is already authenticated. If not, we
		// redirect them to the login endpoint.
		// user := authenticated(r)
		// if user == "" {
		// 	http.Redirect(w, r, "/login?consent="+consentRequestID, http.StatusFound)
		// 	return
		// }

		// Apparently, the user is logged in. Now we check if we received POST
		// request, or a GET request.
		if r.Method == "POST" {
			// Ok, apparently the user gave their consent!

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

			putUrl := fmt.Sprintf("%s/accept?%s", hConf.ConsentRequestRoute, params.Encode())

			body := &BodyAcceptOAuth2Consent{
				GrantScope:               []string{"openid"},
				GrantAccessTokenAudience: []string{"http://localhost:3002"},
				Remember:                 false,
				RememberFor:              3600,
				Session:                  struct{}{},
			}

			putAndRedirect(putUrl, body, w, r, http.DefaultClient)
			return
		}

		// We received a get request, so let's show the html site where the user may give consent.
		fillTemplate := struct {
			ConsentRequestID string
			ClientID         string
			RequestedScopes  []string
		}{
			ConsentRequestID: challenge,
			ClientID:         "app id",
			RequestedScopes:  []string{"scope1", "scope2"},
		}

		renderTemplate(w, "consent.html", fillTemplate)
	}
}
