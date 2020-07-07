package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/common/env"

	"main/provider"
	"main/users"
)

// This store will be used to save user authentication
var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

func main() {

	hConf := &provider.HydraConfig{
		LoginRequestRoute:   "http://localhost:4445/oauth2/auth/requests/login",
		ConsentRequestRoute: "http://localhost:4445/oauth2/auth/requests/consent",
	}
	db := users.NewDB()

	// Context we want handlers to have access to
	envContext := &provider.AuthContext{hConf, db}

	router := httprouter.New()

	// Set up a router and some routes
	router.GET("/login", envContext.GetLogin)
	router.POST("/login", envContext.PostLogin)
	router.GET("/consent", envContext.GetConsent)
	router.POST("/consent", envContext.PostConsent)

	// TODO
	router.POST("/signup", envContext.HandleSignup)

	// Start http server
	log.Println("Listening on :" + env.Getenv("PORT", "3002"))
	log.Fatal(http.ListenAndServe(":"+env.Getenv("PORT", "3002"), router))
}
