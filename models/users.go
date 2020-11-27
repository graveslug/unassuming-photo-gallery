package models

import (
	"errors"

	"github.com/graveslug/unassuming-photo-gallery/hash"
	"github.com/graveslug/unassuming-photo-gallery/rand"

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

	_ UserDB      = &userGorm{}
	_ UserService = &userService{}
)

const hmacSecretKey = "thisWillChangeToo"

//UserDB is used to interact with the users database
//For pretty much all single user queries:
//If the user is found, we will return a nil error
//If the user is not found, we will return ErrNotFound
//If there is another error, we will return an error with more information about what went wrong. This may not be an error generated by the models package.
//For single user queries, any err but ErrNotFound should probaly result in a 500 error until we make "public" facing errors
type UserDB interface {
	//Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//Used to close a DB connection
	Close() error

	//Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

//User model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//UserService is a set of methods user to manipulate and work with the user model
type UserService interface {
	//Authetnticate will verify the provided email address and
	//password are correct. If correct the user is returned to that email
	Authenticate(email, password string) (*User, error)
	UserDB
}

//UserService  Responsible for implementing the Authenticate method only
type userService struct {
	UserDB
}

//userGorm represents our database interaction layer
//and implements the userDB interface fully
type userGorm struct {
	db *gorm.DB
}

//userValidator is the validation layer that validates
//and normalizes data before passing it on to the next
//UserDB in our interface chain
type userValidator struct {
	UserDB
	hmac hash.HMAC
}

//userValFn define a function format. It accepts any function that accepts
//a pointer to a user and returns an error
//ie all validation functions are going to need to accept a User pointer
//do whatever work is needed and return either a nil
//or an error during the validation
type userValFn func(*User) error

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

//NewUserService opens a connection to the database.
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

//ByID will look up the user with the provided ID.
//if user found returns the user if not returns the error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Close the UserService database connection.
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

//DestructiveReset drops the user table and rebuilds it. Not for production build
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

//Create will create the prodvided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields
//Notably when the function first starts it will generate
//a default Remember token before it tries to hash the token.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

//Create will create the user and backfill data like the ID createdAt and updatedAt fields
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
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
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.bcryptPassword,
		uv.hmacRemember)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

//Update updates the user
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

//Delete the user. BUT the thing with Gorm is that if we provide a id of zero to Gorm-Gorm will delete all our users. To prevent this we wrote an error variable and set an if statement to zero.
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

//AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

//ByRemember looks up a user with a given rememberToken
//and returns that user. This method will handle hashing the token for us.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		//We DO NOT need to run this if the password hasn't been changed
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

//runUserValFns accepts a pointer and any number of validation functions
// as its arguments then it iterates over each validation function
//by using a for loop with the range.
//As it iterates over each function, it calls the validation functions
// passing in the user as the argument and catching the error return value.
//If the return value is nil, no error
//we continue to the next validation function.
//If an error is recieved we stop the validations and return that error.
//When all validations successfully run we can return nil,
//which says no error found.
//the pointer on User allows us to persist all changes
//between each function allowing us to capture any data normalization
//that happens along the way
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

//Since we cannot save users to the database without a Remember Token
//we must verify that remember tokens exists and if not
//a value is given by using the rand package we created.
//This only runs when we are creating a new user
//
//While its tempting to validate the length/size of the token: don't.
//It will needlessly complicate the function here and lead to bugs.
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}
