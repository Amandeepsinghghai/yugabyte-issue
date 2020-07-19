package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5433
	user     = "yugabyte"
	password = "yugabyte"
	dbname   = "yugabyte"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	var createStmt = `
		CREATE TABLE IF NOT EXISTS users (
			id int PRIMARY KEY,
			name varchar,
			age int,
			language varchar
		)
	`
	if _, err := db.Exec(createStmt); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created table users")

	// Insert into the table.
	var insertStmt string = "INSERT INTO users(id, name, age, language) VALUES (1, 'John', 35, 'Go') ON CONFLICT DO NOTHING"
	if _, err := db.Exec(insertStmt); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted data: %s\n", insertStmt)

	// Read from the table.
	var name string
	var age int
	var language string

	rows, err := db.Query(`SELECT name, age, language FROM users WHERE id = 1`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("Query for id=1 returned: ")
	for rows.Next() {
		err := rows.Scan(&name, &age, &language)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Row[%s, %d, %s]\n", name, age, language)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`UPDATE employee SET name = 'doe' WHERE id = 1`)
	if err != nil {
		log.Fatal(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated Row Count: %d\n", count)

	defer db.Close()
}
