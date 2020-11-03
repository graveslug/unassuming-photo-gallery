package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "graveslug"
	dbname = "unassuming_photo_gallery"
)

//User a unique model with constraints for input fields of how we want our data to come in
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	// Orders []Order
}

//Order model for orders
type Order struct {
	gorm.Model
	//unsigned integer because we will never use negative IDs
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)
	//Note: AutoMigrate will only create things that don't exist
	db.AutoMigrate(&User{})

	var user User
	db.First(&user)
	if db.Error != nil {
		panic(db.Error)
	}

	createOrder(db, user, 1001, "Fake Description #1")
	createOrder(db, user, 9999, "Fake Description #2")
	createOrder(db, user, 8800, "Fake description #3")
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name, email
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}

// //query the first user
// var u User
// 	db.First(&u)
// 	if db.Error != nil {
// 		panic(db.Error)
// 	}
// 	fmt.Println(u)

// //query user if ID of 1
// var u User
// id := 1
// db.First(&u, id)
// if db.Error != nil {
// 	panic(db.Error)
// }
// fmt.Println(u)

// //Query with a where() method
// var u User
// maxId := 3

// db.Where("id <= ?", maxId).First(&u)
// if db.Error != nil {
// 	panic(db.Error)
// }
// fmt.Println(u)

// //Querying with an existing user
// var u User
// u.Email = "someUser@thisEmail.io"

// db.Where(u).First(&u)
// if db.Error != nil {
// 	panic(db.Error)
// }
// fmt.Println(u)

//query muiltiple records with gorm

// var users []User
// db.Find(&users)
// if db.Error != nil {
// 	panic(db.Error)
// }
// fmy.Println("Retrieved", len(users), "users.")
// fmt.Println(users)
