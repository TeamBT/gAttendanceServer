package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//Student foldy flaps
type Student struct {
	id       int
	name     string
	rfid     string
	password string
	partial  int
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/gAT?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	rows, err := db.Query("SELECT * FROM student;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	studs := make([]Student, 0)
	for rows.Next() {
		stud := Student{}
		err := rows.Scan(&stud.id, &stud.name, &stud.rfid, &stud.password, &stud.partial)
		if err != nil {
			panic(err)
		}
		studs = append(studs, stud)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, stud := range studs {
		// fmt.Println(bk.isbn, stud.title, bk.author, bk.price)
		fmt.Printf("%s, %s, %s, %b\n", stud.name, stud.rfid, stud.password, stud.partial)
	}
}
