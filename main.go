package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

//Student ...
type Student struct {
	id       string
	Name     string
	Rfid     string
	Password string
	Partial  string
	Here     bool
	Excused  bool
}

func main() {
	// os.Setenv("PORT", "8080")
	http.HandleFunc("/", redirectStudent)
	http.HandleFunc("/student", studentsIndex)
	http.HandleFunc("/student/show", studentShow)
	http.HandleFunc("/student/update", studentUpdateProcess)
	http.HandleFunc("/student/reset", resetStudents)
	http.HandleFunc("/student/delete", deleteStudent)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func redirectStudent(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/student", http.StatusSeeOther)
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
		err = rows.Scan(&stud.id, &stud.Name, &stud.Rfid, &stud.Password, &stud.Partial, &stud.Here, &stud.Excused)

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

func studentUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	stud := Student{}
	stud.id = r.FormValue("id")
	stud.Rfid = r.FormValue("rfid")
	stud.Here = r.FormValue("here") == "true"
	stud.Excused = r.FormValue("excused") == "true"

	// insert values
	if stud.id != "" && stud.Rfid == "" {
		_, err := db.Exec("UPDATE student SET here=$2, excused=$3 WHERE id=$1;", stud.id, stud.Here, stud.Excused)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	if stud.Rfid != "" && stud.id == "" {
		_, err := db.Exec("UPDATE student SET here=$2, excused=$3 WHERE rfid=$1;", stud.Rfid, stud.Here, stud.Excused)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}
}

func resetStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := db.Exec("UPDATE student SET here=$1, excused=$2;", false, false)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/student", http.StatusSeeOther)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// delete book
	_, err := db.Exec("DELETE FROM student WHERE id=$1;", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student", http.StatusSeeOther)
}
