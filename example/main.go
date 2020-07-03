package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// This store will be used to save user authentication
var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

var clientConfig = oauth2.Config{
	ClientID:     "open-id-client",
	ClientSecret: "secret",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "http://localhost:4444/oauth2/auth",
		TokenURL: "http://localhost:4444/oauth2/token",
	},
	RedirectURL: "http://localhost:3003/callback",
	Scopes:      []string{"openid"},
}

// A state for performing the OAuth 2.0 flow. This is usually not part of a consent app, but in order for the demo
// to make sense, it performs the OAuth 2.0 authorize code flow.
var state = "demostatedemostatedemo"

func main() {
	var err error

	if err != nil {
		log.Fatalf("Unable to connect to the Hydra SDK because %s", err)
	}

	// Set up a router and some routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/callback", handleCallback)

	// Start http server
	log.Println("Listening on :3003")
	log.Fatal(http.ListenAndServe(":3003", nil))
}

// handles request at /home - a small page that let's you know what you can do in this app. Usually the first.
// page a user sees.
func handleHome(w http.ResponseWriter, _ *http.Request) {

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", "open-id-client")
	params.Add("redirect_uri", "http://localhost:3003/callback")
	params.Add("scope", "openid")
	params.Add("state", state)
	params.Add("aud", "http://localhost:3003")

	req, _ := http.NewRequest("GET", "http://localhost:4444/oauth2/auth?"+params.Encode(), nil)

	renderTemplate(w, "home.html", req.URL)
}

// Once the user has given their consent, we will hit this endpoint. Again,
// this is not something that would be included in a traditional consent app,
// but we added it so you can see the data once the consent flow is done.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	// in the real world you should check the state query parameter, but this is omitted for brevity reasons.

	token, err := clientConfig.Exchange(context.Background(), r.URL.Query().Get("code"))

	if err != nil {
		http.Error(w, errors.Wrap(err, "Could not exhange token").Error(), http.StatusBadRequest)
		return
	}

	// Render the output
	renderTemplate(w, "callback.html", struct {
		*oauth2.Token
		IDToken interface{}
	}{
		Token:   token,
		IDToken: token.Extra("id_token"),
	})
}

// renderTemplate is a convenience helper for rendering templates.
func renderTemplate(w http.ResponseWriter, id string, d interface{}) bool {
	if t, err := template.New(id).ParseFiles("./templates/" + id); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	} else if err := t.Execute(w, d); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	}
	return true
}
