package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
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
			last_logged_at timestamp
		);
	`

	if _, err := db.Exec(createStmt); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created table users")

	// Insert into the table.
	var insertStmt string = "INSERT INTO users(id, last_logged_at) VALUES (1, NOW()) ON CONFLICT DO NOTHING"
	if _, err := db.Exec(insertStmt); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted data: %s\n", insertStmt)

	plainSQL := func(t time.Duration) {
		for i := 0; i < 50; i++ {
			time.Sleep(t)
			_, err = db.Query("Select * FROM users where id = 1")
			if err != nil {
				log.Fatal(err)
			}

			_, err = db.Exec("UPDATE users SET last_logged_at = NOW() WHERE id = 1")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Loop count: %d\n", i)
		}
	}

	go plainSQL(50 * time.Millisecond)
	go plainSQL(50 * time.Millisecond)
	time.Sleep(5000 * time.Millisecond)
	fmt.Println("Printing from main")

	defer db.Close()
}
