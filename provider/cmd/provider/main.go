package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/common/env"

	"provider/internal/users"
)

// This store will be used to save user authentication
var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

type Env struct {
	hConf *hydraConfig
	db    users.UserStore
}

func main() {

	hConf := &hydraConfig{
		LoginRequestRoute:   "http://localhost:4445/oauth2/auth/requests/login",
		ConsentRequestRoute: "http://localhost:4445/oauth2/auth/requests/consent",
	}
	db := users.ConnectDB()

	// Context we want handlers to have access to
	envContext := &Env{hConf, db}

	router := httprouter.New()

	// Set up a router and some routes
	router.GET("/login", envContext.getLogin)
	router.POST("/login", envContext.postLogin)
	router.GET("/consent", envContext.getConsent)
	router.POST("/consent", envContext.postConsent)

	// TODO
	router.POST("/signup", envContext.handleSignup)

	// Start http server
	log.Println("Listening on :" + env.Getenv("PORT", "3002"))
	log.Fatal(http.ListenAndServe(":"+env.Getenv("PORT", "3002"), router))
}
