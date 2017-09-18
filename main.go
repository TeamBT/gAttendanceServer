package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://jwuzvvbakkswqj:638c8277583ccec74041d91560bd319728949fe1a8085eb0570ad9e7e29cbf3c@ec2-54-243-255-57.compute-1.amazonaws.com:5432/derfn0d9nq69b3")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

// export fields to templates
// fields changed to uppercase
type Student struct {
	id       string
	Name     string
	Rfid     string
	Password string
	Partial  string
	Here     string
	Excused  string
}

func main() {
	http.HandleFunc("/student", studentsIndex)
	http.HandleFunc("/student/show", studentShow)
	// http.HandleFunc("/books/create", booksCreateForm)
	// http.HandleFunc("/books/create/process", booksCreateProcess)
	// http.HandleFunc("/student/update", studentUpdateForm)
	http.HandleFunc("/student/update", studentUpdateProcess)
	// http.HandleFunc("/books/delete/process", booksDeleteProcess)
	http.ListenAndServe(":8080", nil)
}

func studentsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	studs := make([]Student, 0)
	for rows.Next() {

		stud := Student{}
		err := rows.Scan(&stud.id, &stud.Name, &stud.Rfid, &stud.Password, &stud.Partial, &stud.Here, &stud.Excused) // order matters

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		studs = append(studs, stud)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	js, err := json.Marshal(studs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func studentShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM student WHERE id = $1", id)

	stud := Student{}
	err := row.Scan(&stud.id, &stud.Name, &stud.Rfid, &stud.Password, &stud.Partial, &stud.Here, &stud.Excused)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(stud)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//
// func booksCreateForm(w http.ResponseWriter, r *http.Request) {
// }
//
// func booksCreateProcess(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	// get form values
// 	bk := Book{}
// 	bk.Isbn = r.FormValue("isbn")
// 	bk.Title = r.FormValue("title")
// 	bk.Author = r.FormValue("author")
// 	p := r.FormValue("price")
//
// 	// validate form values
// 	if bk.Isbn == "" || bk.Title == "" || bk.Author == "" || p == "" {
// 		http.Error(w, http.StatusText(400), http.StatusBadRequest)
// 		return
// 	}
//
// 	// convert form values
// 	f64, err := strconv.ParseFloat(p, 32)
// 	if err != nil {
// 		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the price", http.StatusNotAcceptable)
// 		return
// 	}
// 	bk.Price = float32(f64)
//
// 	// insert values
// 	_, err = db.Exec("INSERT INTO books (isbn, title, author, price) VALUES ($1, $2, $3, $4)", bk.Isbn, bk.Title, bk.Author, bk.Price)
// 	if err != nil {
// 		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 		return
// 	}
// }
//
// func booksUpdateForm(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	isbn := r.FormValue("isbn")
// 	if isbn == "" {
// 		http.Error(w, http.StatusText(400), http.StatusBadRequest)
// 		return
// 	}
//
// 	row := db.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)
//
// 	bk := Book{}
// 	err := row.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		http.NotFound(w, r)
// 		return
// 	case err != nil:
// 		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 		return
// 	}
// }
//
func studentUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	stud := Student{}
	stud.id = r.FormValue("id")
	// stud.Rfid = r.FormValue("rfid")
	stud.Here = r.FormValue("here")
	stud.Excused = r.FormValue("excused")

	// insert values
	_, err := db.Exec("UPDATE student SET id = $1, here=$2, excused=$3 WHERE id=$1;", stud.id, stud.Here, stud.Excused)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

}

//
// func booksDeleteProcess(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	isbn := r.FormValue("isbn")
// 	if isbn == "" {
// 		http.Error(w, http.StatusText(400), http.StatusBadRequest)
// 		return
// 	}
//
// 	// delete book
// 	_, err := db.Exec("DELETE FROM books WHERE isbn=$1;", isbn)
// 	if err != nil {
// 		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
// 		return
// 	}
//
// 	http.Redirect(w, r, "/books", http.StatusSeeOther)
// }
