package main

import
(
  "fmt"
  "github.com/codegangsta/martini"
  "github.com/russross/blackfriday"
  "net/http"
  _ "github.com/lib/pq"
  "database/sql"
)

func SetupDB() *sql.DB {
  db, err := sql.Open("postgres", "lesson4 sslmode=disable")
  PanicIf(err)
  return db
}

func PanicIf(err error)  {
  if err != nil  {
    panic(err)
  }
}
func main()  {
  m := martini.Classic()
  m.Map(SetupDB())

  //m.Get("/", func(rw http.ResponseWriter, r *http.Request) {
  //  rw.Write([]byte("Hello World"))
  //})

  m.Get("/", func(rw http.ResponseWriter, r *http.Request, db *sql.DB) {
    searchTerm := "%" + r.URL.Query().Get("search") + "%"
    rows, err := db.Query(`SELECT title, author, description FROM  books 
                                 WHERE title ILIKE $1
                                 OR author ILIKE $1
                                 OR description ILIKE $1`, searchTerm)
    PanicIf(err)
    defer rows.Close()

    var title, author, description string
    for rows.Next() {
      err := rows.Scan(&title, &author, &description)
      PanicIf(err)
      fmt.Fprintf(rw, "Title: %s\nAuthor: %s\n Description: %s\n\n", title, author, description)
    }
  })

  m.Post("/generate", func(r *http.Request) []byte {
    body := r.FormValue("body")
    return blackfriday.MarkdownBasic([]byte(body))
  })


  m.Run()
}