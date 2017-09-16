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
	here     bool
	excused  bool
}

func main() {
	db, err := sql.Open("postgres", "postgres://jwuzvvbakkswqj:638c8277583ccec74041d91560bd319728949fe1a8085eb0570ad9e7e29cbf3c@ec2-54-243-255-57.compute-1.amazonaws.com:5432/derfn0d9nq69b3")
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
		err := rows.Scan(&stud.id, &stud.name, &stud.rfid, &stud.password, &stud.partial, &stud.here, &stud.excused)
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
		fmt.Printf("%d, %s, %s, %s, %b, %t, %t\n", stud.id, stud.name, stud.rfid, stud.password, stud.partial, stud.here, stud.excused)
	}
}
