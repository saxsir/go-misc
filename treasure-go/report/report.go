package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Report struct {
	Title string
	Body  string
}

func (r *Report) Insert(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare(`insert into reports (title, body) values (?, ?)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(r.Title, r.Body)
}

func Select(db *sql.DB) ([]Report, error) {
	rows, err := db.Query(`select todo_id, title, body from reports`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []Report{}
	for rows.Next() {
		var title string
		var body string
		if err := rows.Scan(&title, &body); err != nil {
			return nil, err
		}
		r := Report{Title: title, Body: body}
		reports = append(reports, r)
	}
	return reports, nil
}

func main() {
	var (
		title = flag.String("title", "kuke", "cool title")
		body  = flag.String("body", "treasure", "cool body")
	)
	flag.Parse()
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/treasure")
	if err != nil {
		log.Fatalf("open failed: %s", err)
	}
	defer db.Close()

	r := &Report{*title, *body}
	if _, err := r.Insert(db); err != nil {
		log.Fatalf("insert failed: %s", err)
	}
	log.Printf("report added: %#v\n", r)

}
