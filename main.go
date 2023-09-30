package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

type StatementStruct struct {
	Id             int
	Name           string
	LastName       string
	Date           string
	Status         int
	Statement      string
	PassportSeries string
	Time           time.Time
}

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
}

var users = []User{}

func AddStatement(page http.ResponseWriter, r *http.Request) {

	if len(users) > 0 {
		tmpl, err := template.ParseFiles("html_files/addstatement.html", "html_files/header.html")
		if err != nil {
			panic(err)
		}

		tmpl.ExecuteTemplate(page, "addstatement", users[0])
	} else {
		tmpl, err := template.ParseFiles("html_files/warning.html", "html_files/header.html")
		if err != nil {
			panic(err)
		}
		tmpl.ExecuteTemplate(page, "warning", nil)
	}

}

func AddStatementPost(page http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	lastname := r.FormValue("lastname")
	date := r.FormValue("date")
	statement := r.FormValue("statement")
	passportseries := r.FormValue("passportseries")
	userid := r.FormValue("id")
	time := time.Now()

	connStr := "user=postgres password=123456 dbname=mygovdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO public.statements (name, lastname, date, status, statement, passportseries, time, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", name, lastname, date, 100, statement, passportseries, time, userid)

	http.Redirect(page, r, "/", http.StatusSeeOther)
}

func index(page http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	deserializedUser := User{}
	err = json.Unmarshal([]byte(string(b)), &deserializedUser)

	users = append(users, deserializedUser)
	fmt.Println(deserializedUser.Token)
	http.Redirect(page, r, "/", http.StatusSeeOther)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", AddStatement)
	http.HandleFunc("/adding_statement", AddStatementPost)
	http.HandleFunc("/index", index)

	http.ListenAndServe(":8080", nil)
}
