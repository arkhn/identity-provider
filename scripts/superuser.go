package main

import (
	"fmt"
	"main/users"
	"os"
)

func main() {
	databaseHost := os.Getenv("PROVIDER_DB_HOST")
	port := os.Getenv("PROVIDER_DB_PORT")
	username := os.Getenv("PROVIDER_DB_USER")
	password := os.Getenv("PROVIDER_DB_PASSWORD")
	databaseName := os.Getenv("PROVIDER_DB_NAME")

	superuserEmail, ok := os.LookupEnv("SUPERUSER_EMAIL")
	if !ok {
		superuserEmail = "admin@arkhn.com"
	}
	superuserPassword, ok := os.LookupEnv("SUPERUSER_PASSWORD")
	if !ok {
		panic("SUPERUSER_PASSWORD env variable is required")
	}

	db, err := users.NewDB(databaseHost, port, username, password, databaseName)
	if err != nil {
		panic(err)
	}

	user := &users.User{
		Name:     "admin",
		Email:    superuserEmail,
		Password: superuserPassword,
	}

	createdUser, err := db.AddUser(user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created user %s with email %s\n", createdUser.Name, createdUser.Email)
	os.Exit(0)
}
