package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Amandeepsinghghai/yugabyte-issue/models"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	plainSQL := func() {
		for i := 0; i < 50; i++ {
			time.Sleep(25 * time.Millisecond)
			_, err = db.Exec("UPDATE users SET last_logged_at = NOW() WHERE id = 1")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Update loop count: %d\n", i)
		}
	}

	ctx := context.Background()

	// SQLBoiler
	sqlBoiler := func() {
		for i := 0; i < 50; i++ {
			time.Sleep(50 * time.Millisecond)
			user, err := models.FindUser(ctx, db, 1)
			if err != nil {
				log.Fatal(err)
			}

			_, err = user.Update(ctx, db, boil.Whitelist("last_logged_at"))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Update loop count sqlboiler: %d\n", i)
		}
	}

	go plainSQL()
	go sqlBoiler()
	time.Sleep(10000 * time.Millisecond)
	fmt.Println("Printing from main")

	defer db.Close()
}
