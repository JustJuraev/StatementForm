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
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func AddStatement(page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html_files/addstatement.html", "html_files/header.html")
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	deserializedUser := User{}
	err = json.Unmarshal([]byte(string(b)), &deserializedUser)
	//fmt.Println(deserializedUser.Token)

	//client := redis.NewClient(&redis.Options{
	//		Addr:     "localhost:2222",
	//		Password: "", // no password set
	//		DB:       0,  // use default DB
	//})
	//	ctx := context.Background()
	//var userSession = map[string]string{}
	//userSession = client.HGetAll(ctx, "user-session:123").Val()
	fmt.Println(string(b))
	fmt.Println(deserializedUser.Token)
	if string(b) != "" {
		fmt.Println(32131)
		livingTime := 1800 * time.Hour
		expiration := time.Now().Add(livingTime)
		cookie := http.Cookie{Name: "token", Value: deserializedUser.Token, Expires: expiration}
		http.SetCookie(page, &cookie)
		return
	}

	//fmt.Println(userSession)

	tmpl.ExecuteTemplate(page, "addstatement", nil)
}

func AddStatementPost(page http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	lastname := r.FormValue("lastname")
	date := r.FormValue("date")
	statement := r.FormValue("statement")
	passportseries := r.FormValue("passportseries")
	time := time.Now()

	connStr := "user=postgres password=123456 dbname=mygovdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO public.statements (name, lastname, date, status, statement, passportseries, time) VALUES ($1, $2, $3, $4, $5, $6, $7)", name, lastname, date, 100, statement, passportseries, time)

	http.Redirect(page, r, "/", http.StatusSeeOther)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", AddStatement)
	http.HandleFunc("/adding_statement", AddStatementPost)

	http.ListenAndServe(":8080", nil)
}
