package users

import (
	"errors"
	"fmt"

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

func NewDB(databaseHost string, port string, username string, password string, databaseName string) (*DB, error) {

	dbURI := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		databaseHost, port, username, databaseName, password,
	)

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		return nil, err
	}

	// Migrate the schema
	// TODO check doc about that
	db.AutoMigrate(&User{})

	return &DB{db}, nil
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
