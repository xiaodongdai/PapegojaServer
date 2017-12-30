package main

import (
    "fmt"
    "net/http"
    "flag"
     "html/template"
     "database/sql"
    _ "github.com/lib/pq"
)

const (
    DB_HOST     = "localhost"
    DB_PORT     = 5431
    DB_USER     = "postgres"
    DB_PASSWORD = "example"
    DB_NAME     = "postgres"
)

type User struct {
  userName string
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func handler2(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there 2 I love %s!", r.URL.Path[1:])
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}



func writeDB(user string, email string) {
  dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
      DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
  fmt.Println("dbinfo:" + dbinfo)
  db, err := sql.Open("postgres", dbinfo)
  checkErr(err)
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("# Inserting values :" + user)
  user = "aaa"
  email = "bbb"
  var lastInsertId int
  err = db.QueryRow("INSERT INTO test(user_name,email_address) VALUES($1,$2) returning id;", "aaa", "bbb").Scan(&lastInsertId)
  checkErr(err)
}



func handlerLogin(w http.ResponseWriter, r *http.Request) {
  // temporarily use the HTTP basic authentication, 
  // receive username and password in base64 format.
  username, _, e := r.BasicAuth()
  if e == false {
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprintf(w, "403 - Fail!")
    return
  }

  dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
      DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
  db, err := sql.Open("postgres", dbinfo)
  checkErr(err)
  defer db.Close()
  // check the user database
  rows, err := db.Query("SELECT count(*) FROM userinfo WHERE user_name = $1", username)
  checkErr(err)
  defer rows.Close()
  // receive the count.
  var count int  
  for rows.Next() {
    err = rows.Scan(&count)
    checkErr(err)
    fmt.Printf("%3v", count)
  }
  if count == 1 {
    fmt.Fprintf(w, "200 - Succeed!")
  } else {
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprintf(w, "403 - Failed!")
  }
}


func handleActive(w http.ResponseWriter, r *http.Request) {
  fmt.Println("path:" + r.URL.Path[1:])
  user := r.URL.Path[1]
  email :=  r.URL.Path[2]
  writeDB(string(user), string(email))
  fmt.Fprintf(w, "Success %s!", r.URL.Path[0:])
}

func handleResetPassword(w http.ResponseWriter, r *http.Request) {
  //fmt.Fprintf(w, "Hi registry I love %s!", r.URL.Path[0:])
  //user := r.URL.Path[1]
  t := template.Must(template.ParseFiles("./templates/resetPassword.html"))
  t.Execute(w, "TEST")
}

func main() {
    directory := flag.String("d", "./images/", "the directory of static file to host")
    flag.Parse()
    http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(*directory))))
    //http.Handle("/images", handler)
    fmt.Println(*directory)
    http.HandleFunc("/", handler)
    http.HandleFunc("/login", handlerLogin)
    http.HandleFunc("/active", handleActive)
    http.HandleFunc("/resetPassword", handleResetPassword)
    http.HandleFunc("/test1", handler2)

    http.ListenAndServe(":8088", nil)
}
