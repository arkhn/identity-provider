package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// The user hits this endpoint if not authenticated. In this example, they can sign in with the credentials
// buzz:lightyear
func handleLogin(hConf *hydraConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		challenge := r.URL.Query().Get("login_challenge")

		// Get info about current flow
		params := url.Values{}
		params.Add("login_challenge", challenge)

		getUrl := fmt.Sprintf("%s?%s", hConf.LoginRequestRoute, params.Encode())
		resp, _ := http.Get(getUrl)

		b, _ := ioutil.ReadAll(resp.Body)
		var jsonResp interface{}
		json.Unmarshal(b, &jsonResp)
		log.Println(jsonResp)
		// TODO do stuff with response

		// Is it a POST request?
		if r.Method == "POST" {
			// Parse the form
			if err := r.ParseForm(); err != nil {
				http.Error(w, errors.Wrap(err, "Could not parse form").Error(), http.StatusBadRequest)
				return
			}

			// Check the user's credentials
			if r.Form.Get("username") != "a" || r.Form.Get("password") != "a" {
				http.Error(w, "Provided credentials are wrong, try buzz:lightyear", http.StatusBadRequest)
				return
			}

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
			putUrl := fmt.Sprintf("%s/accept?%s", hConf.LoginRequestRoute, params.Encode())

			// TODO properly fill body
			body := &BodyAcceptOAuth2Login{
				Subject:     "bob",
				Remember:    false,
				RememberFor: 3600,
				Acr:         "..",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest("PUT", putUrl, bytes.NewBuffer(jsonBody))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				panic(err.Error())
			}

			jsonResp := RedirectResp{}
			json.Unmarshal(b, &jsonResp)

			http.Redirect(w, r, jsonResp.RedirectTo, http.StatusFound)
			return
		}

		// It's a get request, so let's render the template
		renderTemplate(w, "login.html", challenge)
	}
}
