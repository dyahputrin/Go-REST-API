//what to improve ? id use UUID will be better

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tkpdmput24"
	dbname   = "database1"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//insert("Will Smith", "willsmith@gmail.com", psqlInfo)

	//update(6, "Billy Joel", "billy@gmail.com", psqlInfo)

	//delete(6, psqlInfo)

	// err = db.Ping()
	// if err != nil {
	//  panic(err)
	// }

	//fmt.Println(uuid.NewString())

	fmt.Println("Rest API v2.0 - Mux Routers")

	handleRequests()

}

func insert(id string, name string, email string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   INSERT INTO users (id, name, email)
   VALUES ($1, $2, $3)
   RETURNING id`
	id = uuid.NewString() //pake uuid soalnya pasti unique

	err = db.QueryRow(sqlStatement, id, name, email).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}

func update(id string, newName string, newEmail string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   UPDATE users
   SET name = $2, email = $3
   WHERE id = $1;`

	idInt, err := strconv.Atoi(id)

	_, err = db.Exec(sqlStatement, idInt, newName, newEmail)
	if err != nil {
		panic(err)
	}

	fmt.Println("Update record ID:", id)
}

func delete(id string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   DELETE FROM users
   WHERE id = $1;`

	idInt, err := strconv.Atoi(id)

	_, err = db.Exec(sqlStatement, idInt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Delete record ID:", idInt)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	//fmt.Println("Endpoint Hit: homePage")
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", returnAllUsers)
	myRouter.HandleFunc("/user/{id}", returnSingleUser)
	myRouter.HandleFunc("/user", createNewUser).Methods("POST")
	myRouter.HandleFunc("/user/delete/{id}", deleteUser)
	myRouter.HandleFunc("/user/update/{id}", updateUser).Methods("PUT")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var Users []User

func returnAllUsers(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(Users)
	fmt.Println("Endpoint Hit: returnAllUsers")
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	//keyInt, err := strconv.Atoi(key)

	//fmt.Fprintf(w, "Key: "+key)
	for _, user := range Users {
		if user.Id == key {
			json.NewEncoder(w).Encode(user)
		}
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)

	Users = append(Users, user)

	json.NewEncoder(w).Encode(user)

	insert(user.Id, user.Name, user.Email, psqlInfo)

	fmt.Println("Create New User Success")

}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	vars := mux.Vars(r)

	id := vars["id"]

	var user User

	user.Id = id

	for index, user := range Users {
		if user.Id == id {
			Users = append(Users[:index], Users[index+1:]...)
		}
	}

	delete(id, psqlInfo)

	fmt.Println("Delete User Success")

}

func updateUser(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	vars := mux.Vars(r)
	id := vars["id"]

	var updatedEvent User
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updatedEvent)
	for i, user := range Users {
		if user.Id == id {

			user.Id = updatedEvent.Id
			user.Name = updatedEvent.Name
			user.Email = updatedEvent.Email
			Users[i] = user
			json.NewEncoder(w).Encode(user)
		}
	}

	update(updatedEvent.Id, updatedEvent.Name, updatedEvent.Email, psqlInfo)

	fmt.Println("Update User Success")
}
