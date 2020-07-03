package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ory/common/env"
)

// This store will be used to save user authentication
var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

// A state for performing the OAuth 2.0 flow. This is usually not part of a consent app, but in order for the demo
// to make sense, it performs the OAuth 2.0 authorize code flow.
var state = "demostatedemostatedemo"

func main() {

	hConf := hydraConfig{
		LoginRequestRoute:   "http://localhost:4445/oauth2/auth/requests/login",
		ConsentRequestRoute: "http://localhost:4445/oauth2/auth/requests/consent",
	}

	// Set up a router and some routes
	http.HandleFunc("/login", handleLogin(&hConf))
	http.HandleFunc("/consent", handleConsent(&hConf))

	// Start http server
	log.Println("Listening on :" + env.Getenv("PORT", "3002"))
	log.Fatal(http.ListenAndServe(":"+env.Getenv("PORT", "3002"), nil))
}
