package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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
	// 	panic(err)
	// }

	//fmt.Println(uuid.NewString())

}

func insert(id string, name string, email string, psqlInfo string) {
	db, err := sql.Open("postgres", psqlInfo)

	sqlStatement := `
	INSERT INTO users (name, email)
	VALUES ($1, $2)
	RETURNING id`
	id = uuid.NewString() //we use uuid to make sure it's unique

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
