package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/common/env"

	"main/provider"
	"main/users"
)

// This store will be used to save user authentication
// var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

func main() {

	hConf := &provider.HydraConfig{
		LoginRequestRoute:   "http://localhost:4445/oauth2/auth/requests/login",
		ConsentRequestRoute: "http://localhost:4445/oauth2/auth/requests/consent",
	}
	db, err := users.NewDB()

	if err != nil {
		log.Fatal(err.Error())
	}

	// Context we want handlers to have access to
	envContext := &provider.Provider{
		HConf: hConf,
		Db:    db,
	}

	router := httprouter.New()

	// Set up a router and some routes
	router.GET("/login", envContext.GetLogin)
	router.POST("/login", envContext.PostLogin)
	router.GET("/consent", envContext.GetConsent)
	router.POST("/consent", envContext.PostConsent)

	// TODO
	router.POST("/signup", envContext.HandleSignup)

	// Start http server
	serverUrl := fmt.Sprintf("localhost:%s", env.Getenv("PORT", "3002"))
	fmt.Printf("Listening on: %s\n", serverUrl)
	log.Fatal(http.ListenAndServe(serverUrl, router))
}
