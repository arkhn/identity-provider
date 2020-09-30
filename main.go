package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"main/provider"
	"main/users"
)

// This store will be used to save user authentication
// var store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))

// The session is a unique session identifier
const sessionName = "authentication"

func main() {

	loginRequestRoute := os.Getenv("LOGIN_REQUEST_ROUTE")
	consentRequestRoute := os.Getenv("CONSENT_REQUEST_ROUTE")

	databaseHost := os.Getenv("PROVIDER_DB_HOST")
	databasePort := os.Getenv("PROVIDER_DB_PORT")
	databaseUsername := os.Getenv("PROVIDER_DB_USER")
	databasePassword := os.Getenv("PROVIDER_DB_PASSWORD")
	databaseName := os.Getenv("PROVIDER_DB_NAME")

	hConf := &provider.HydraConfig{
		LoginRequestRoute:   loginRequestRoute,
		ConsentRequestRoute: consentRequestRoute,
	}
	db, err := users.NewDB(databaseHost, databasePort, databaseUsername, databasePassword, databaseName)

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

	router.GET("/signup", envContext.GetSignup)
	router.POST("/signup", envContext.PostSignup)

	// Start http server
	serverUrl := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	fmt.Printf("Listening on: %s\n", serverUrl)
	log.Fatal(http.ListenAndServe(serverUrl, router))
}
