package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id                    int
	email, user, password string
}

var db, dbErr = sql.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/askus")

func fetchData() {
	results, err := db.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	var users []User

	for results.Next() {
		var user User
		err = results.Scan(&user.id, &user.email, &user.user, &user.password)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		users = append(users, user)
	}

	fmt.Println(users)
}

func executeQuery() {
	insert, err := db.Query("INSERT INTO users(email, user, password) VALUES ('redis1@gmail.com', 'redis1', 'redis1')")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

func fetchSingle() {
	var user User

	err := db.QueryRow("SELECT * FROM users where id = ?", 2).Scan(&user.id, &user.user, &user.email, &user.password)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(user)
}

func main() {
	sayHello()

	if dbErr != nil {
		panic(dbErr.Error())
	}

	//fetchSingle()
	fetchData()
	//query()

	defer db.Close()
}
