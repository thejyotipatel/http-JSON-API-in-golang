package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for users
type User struct {
	Id          string `json:"id"`
	FirstName   string `json:"fastname"`
	LastName    string `json:"lastname"`
	DateOfBirth string `json:"dob"`
	Email       string `json:"email"`
	PhoneNumber int    `json:"phone"`
}

// Init user
var users []User

// middleware
func (u *User) IsEmpty() bool {
	return u.Email == "" && u.FirstName == ""
}

// check email patten
func (u *User) checkEmailPatten() bool {
	regText := `/\b[\w\.-]+@[\w\.-]+\.\w{2,4}\b/gi`
	return u.Email != regText
}

func (u *User) checkPhoneNumber() bool {
	regText := `[(]?\d{3}[)]?\s?-?\s?\d{3}\s?-?\s?\d{4}`

	number, _ := strconv.Atoi(regText)

	return u.PhoneNumber == number
}

// get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// get one user
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant-Type", "application/json")

	param := mux.Vars(r) // Get param
	// loop through users
	for _, item := range users {
		if item.Id == param["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

// create user
func ctreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant-Type", "application/json")

	// if parem is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// if email and name empty
	if user.IsEmpty() {
		json.NewEncoder(w).Encode("Please provide name and email")
		return
	}

	rand.Seed(time.Now().UnixNano())
	user.Id = strconv.Itoa(rand.Intn(100))
	users = append(users, user)
	json.NewEncoder(w).Encode(user)

}

// update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant-Type", "application/json")

	param := mux.Vars(r) // Get id of user

	// loop through users
	for index, item := range users {
		if item.Id == param["id"] {
			users = append(users[:index], users[index+1:]...)

			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.Id = param["id"]
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}

	}
}

// delete user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Delete one user")

	w.Header().Set("Contant-Type", "application/json")

	// get id of user
	param := mux.Vars(r)

	// loop throught all user
	for index, item := range users {
		if item.Id == param["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
}

func serverPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome this side and server running at 5000.</h1>"))
}

func main() {
	fmt.Printf("Server is running to port no 5000...")

	// init Router
	r := mux.NewRouter()

	// Mock data
	users = append(users, User{Id: "1", FirstName: "joan", LastName: "smit", DateOfBirth: "12-3-20", Email: "joan@gmail.com", PhoneNumber: 768756789})
	users = append(users, User{Id: "2", FirstName: "max", LastName: "may", DateOfBirth: "12-8-20", Email: "max@gmail.com", PhoneNumber: 435798776})

	// Routing Handlers / Endpoints
	r.HandleFunc("/", serverPage).Methods("GET")
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/api/users", ctreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", r))
}
