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
	superuserPassword := os.Getenv("SUPERUSER_PASSWORD")

	if superuserPassword == "" {
		panic("SUPERUSER_PASSWORD env variable is required")
	}

	db, err := users.NewDB(databaseHost, port, username, password, databaseName)
	if err != nil {
		panic(err)
	}

	user := &users.User{
		Name:     "admin",
		Email:    "admin@arkhn.com",
		Password: superuserPassword,
	}

	createdUser, err := db.AddUser(user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created user %s with email %s\n", createdUser.Name, createdUser.Email)
}
