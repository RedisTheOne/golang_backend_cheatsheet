package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}

	var users []User
	for results.Next() {
		var user User
		err = results.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func checkIfUserExistsInDB(username string, password string) bool {
	var user User
	err := db.QueryRow("SELECT * FROM users where username = '"+username+"' AND password = '"+password+"'").Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return false
	}

	return true
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginRequestBody LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&loginRequestBody)
	if len(loginRequestBody.Password) != 0 && len(loginRequestBody.Username) != 0 {
		if checkIfUserExistsInDB(loginRequestBody.Username, loginRequestBody.Password) {
			response := Response{Message: "Welcome back " + loginRequestBody.Username + "!", Status: true}
			json.NewEncoder(w).Encode(response)
		} else {
			response := Response{Message: "User is not valid.", Status: false}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := Response{Message: "Please fill required fields.", Status: false}
		json.NewEncoder(w).Encode(response)
	}
}

func signUpRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var signupRequestBody SignUpRequest
	_ = json.NewDecoder(r.Body).Decode(&signupRequestBody)

	if len(signupRequestBody.Password) != 0 && len(signupRequestBody.Username) != 0 && len(signupRequestBody.Email) != 0 {
		if checkIfUserExistsInDB(signupRequestBody.Username, signupRequestBody.Password) {
			response := Response{Message: "Username is already taken.", Status: false}
			json.NewEncoder(w).Encode(response)
		} else {
			createUserInDb(signupRequestBody.Username, signupRequestBody.Email, signupRequestBody.Password)
			response := Response{Message: "User created.", Status: true}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := Response{Message: "Please fill required fields.", Status: false}
		json.NewEncoder(w).Encode(response)
	}
}

func createUserInDb(username string, email string, pass string) {
	insert, err := db.Query("INSERT INTO users(email, username, password) VALUES ('" + email + "', '" + username + "', '" + pass + "')")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api", rootRoute).Methods("GET")
	r.HandleFunc("/api/login", loginRoute).Methods("POST")
	r.HandleFunc("/api/signup", signUpRoute).Methods("POST")

	fmt.Println("Server started at port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
