package users

// TODO should we move user management outside the main package?

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	AddUser(user *User) (*User, error)
	FindUser(email string, password string) (*User, error)
}

// TODO What else do we need?
// Authorizations for each client?
// Personal info?
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index"`
	Name     string `gorm:"type:varchar(100)"`
	Password string `json:"Password"`
}

type DB struct {
	*gorm.DB
}

func NewDB() *DB {
	databaseHost := os.Getenv("PROVIDER_DB_HOST")
	port := os.Getenv("PROVIDER_DB_PORT")
	username := os.Getenv("PROVIDER_DB_USER")
	password := os.Getenv("PROVIDER_DB_PASSWORD")
	databaseName := os.Getenv("PROVIDER_DB_NAME")

	dbURI := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		databaseHost, port, username, databaseName, password,
	)

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	// Migrate the schema
	// TODO check doc about that
	db.AutoMigrate(&User{})

	return &DB{db}
}

func (db *DB) FindUser(email string, password string) (*User, error) {
	user := &User{}

	if err := db.Where(&User{Email: email}).First(user).Error; err != nil {
		return nil, errors.New("Email address not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("Invalid login credentials. Please try again")
	}

	return user, nil
}

func (db *DB) AddUser(user *User) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	createdUser := db.Create(user)

	if createdUser.Error != nil {
		return nil, createdUser.Error
	}

	return user, nil
}
