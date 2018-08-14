package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/treasure")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	input := "hogehoge or 1; -- "
	rows, err := db.Query(fmt.Sprintf("select todo_id, title from todos where title = %s", input))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id    int
			title string
		)
		if err := rows.Scan(&id, &title); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d: %s\n", id, title)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
