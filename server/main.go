package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./superdupercloud.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Exec(`CREATE TABLE test (
		id INTEGER PRIMARY KEY,
		first_name TEXT NOT NULL
	)`)
	fmt.Println(res, err)
	res, err = db.Exec("INSERT INTO test VALUES(1, 'abc')")
	fmt.Println(res, err)

	fmt.Println("starting server")
}
