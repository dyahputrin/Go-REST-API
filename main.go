//what to improve ? id use UUID will be better

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	Users = []User{
		User{Id: "1", Name: "Article Description", Email: "Article Content"},
		User{Id: "2", Name: "Article Description", Email: "Article Content"},
	}
	handleRequests()

}

func insert(id string, name string, email string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   INSERT INTO users (name, email)
   VALUES ($1, $2)
   RETURNING id`
	id = uuid.NewString() //pake uuid soalnya pasti unique

	err = db.QueryRow(sqlStatement, name, email).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}

func update(id int, newName string, newEmail string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   UPDATE users
   SET name = $2, email = $3
   WHERE id = $1;`

	_, err = db.Exec(sqlStatement, id, newName, newEmail)
	if err != nil {
		panic(err)
	}

	fmt.Println("Update record ID:", id)
}

func delete(id int, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
   DELETE FROM users
   WHERE id = $1;`

	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Delete record ID:", id)
}

//register nama email (POST); edit id nama/email (PUT); delete id (POST)
//view data nya

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// http.HandleFunc("/", homePage)
	// http.HandleFunc("/users", returnAllUsers)
	// log.Fatal(http.ListenAndServe(":10000", nil))

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", returnAllUsers)
	myRouter.HandleFunc("/user/{id}", returnSingleUser)
	myRouter.HandleFunc("/user", createNewUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Users []User

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllUsers")
	json.NewEncoder(w).Encode(Users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	//fmt.Fprintf(w, "Key: "+key)
	for _, user := range Users {
		if user.Id == key {
			json.NewEncoder(w).Encode(user)
		}
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body

	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)

	Users = append(Users, user)

	json.NewEncoder(w).Encode(user)

	fmt.Println("Create User Success")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	// we then need to loop through all our articles
	for index, user := range Users {
		if user.Id == id {
			Users = append(Users[:index], Users[index+1:]...)
		}
	}
	fmt.Println("Delete User Success")

}

func updateUser(w http.ResponseWriter, r *http.Request) {
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
}
