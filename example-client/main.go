package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var clientConfig = oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:  os.Getenv("AUTH_URL"),
		TokenURL: os.Getenv("TOKEN_URL"),
	},
	RedirectURL: os.Getenv("REDIRECT_URL"),
	Scopes:      []string{"openid", "offline_access"},
}

var state = "demostatedemostatedemo"

func main() {
	// Set up a router and some routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/callback", handleCallback)

	// Start http server
	log.Println("Listening on :3003")
	log.Fatal(http.ListenAndServe(":3003", nil))
}

// handles request at /home - a small page that let's you know what you can do in this app. Usually the first.
// page a user sees.
func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html", nil)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	url := clientConfig.AuthCodeURL(state)
	log.Println(url)
	http.Redirect(w, r, url, http.StatusFound)
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
	if t, err := template.New(id).ParseFiles("example-client/templates/" + id); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	} else if err := t.Execute(w, d); err != nil {
		http.Error(w, errors.Wrap(err, "Could not render template").Error(), http.StatusInternalServerError)
		return false
	}
	return true
}
