package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	//I don't want to defend why I'm doing this. :(
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNotFound makes an apperance when a resource cannot be found in the database. You can keep trying to find it though if you'd like.
	ErrNotFound = errors.New("models: resource not found")
	//ErrInvalidID is returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
	//ErrInvalidPassword is returned when invalid password is used when attempting to authenticate the user
	ErrInvalidPassword = errors.New("models: incorrect password provided")

	userPwPepper = "Don'tGetExcitedThisWillChangeAndThereIsNoCloud"
)

//User model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

//UserService provides methods for querying, creating, and updating the users.
type UserService struct {
	db *gorm.DB
}

//NewUserService opens a connection to the database.
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

// Close the UserService database connection.
func (us *UserService) Close() error {
	return us.db.Close()
}

//ByID will look up the user with the provided ID.
//if user found returns the user if not returns the error
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//DestructiveReset drops the user table and rebuilds it. Not for production build
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

//Create will create the user and backfill data like the ID createdAt and updatedAt fields
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

//first will query using the gorm.DB and it will get thef irst item returned and place it into the dst. If nothing is found in the query it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

//ByEmail finding the user by email search
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//Update updates the user
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//Delete the user. BUT the thing with Gorm is that if we provide a id of zero to Gorm-Gorm will delete all our users. To prevent this we wrote an error variable and set an if statement to zero.
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

//AutoMigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

//Authenticate can be used to authenticate a suer with the provided email address and password
//if the email address provided is invalid, this will return
//nil, ErrNotFound
//If password provided is invalid, this will return
//nil, errInvalidPassword
//if the email and password are both valid this will return
//user, nil
//otherwise if another error is encountered this will return
//nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}
