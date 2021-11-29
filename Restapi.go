package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	_"github.com/lib/pq"
)
const(
	DB_USER = "postgres"
	DB_PASSWORD = "admin"
	DB_NAME = "postgres"
)
//DB setup
func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)

    checkErr(err)

    return db
}

//Book Struct (MODEL) C:\Users\Swift Alliance Ltd\go\src\github.com\gorilla
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	TITLE  string  `json:"title"`
	AFirstname string `json:"afirstname"`
	ALastname string `json:"alastname"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Book `json:"data"`
    Message string `json:"message"`
}

//Init books var as a slice book struct
//var books []Book

// Function for handling messages
func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}

//Check ERROR
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

//get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(books)
	db := setupDB()

    printMessage("Getting books...")

    // Get all book properties from book table
    rows, err := db.Query("SELECT * FROM books Order by id;" )

    // check errors
    checkErr(err)

    // var response []JsonResponse
     var books []Book

    // Foreach books
    for rows.Next() {
        var id int
        var isbn string
        var title string
		var afirstname string
		var alastname string

        err = rows.Scan(&id, &isbn, &title, &afirstname, &alastname)

        // check errors
        checkErr(err)

        books = append(books, Book{ID: strconv.Itoa(id), ISBN: isbn, TITLE: title, AFirstname: afirstname, ALastname : alastname})
    }

    var response = JsonResponse{
    	Type:    "success",
    	Data:    books,
    }

    json.NewEncoder(w).Encode(response)
}

//get single Books
 func getBook(w http.ResponseWriter, r *http.Request) {

  params := mux.Vars(r)
var books []Book
  ID := params["id"]
db := setupDB()
sqlStatement := `SELECT id, isbn, title, afirstname, alastname FROM books WHERE id=$1;`
var id int
var isbn int
var title string
var afirstname string
var alastname string
// Replace 3 with an ID from your database or another random
// value to test the no rows use case.
row := db.QueryRow(sqlStatement, ID)
switch err := row.Scan(&id, &isbn, &title, &afirstname, &alastname); err {
case sql.ErrNoRows:
  fmt.Println("No rows were returned!")
  var response = JsonResponse{
	Type:    "failed",
	Message: "Book With Id number "+ID+" does not exist",
}

json.NewEncoder(w).Encode(response)
case nil:
  fmt.Println(id, isbn, title, afirstname, alastname)
  books = append(books, Book{ID: strconv.Itoa(id), ISBN: strconv.Itoa(isbn), TITLE: title, AFirstname: afirstname, ALastname : alastname})

var response = JsonResponse{
	Type:    "success",
	Data:    books,
}

json.NewEncoder(w).Encode(response)

default:
  panic(err)
}
}

func createBook(w http.ResponseWriter, r *http.Request){
//response and request handlers
bookId := r.FormValue("id")
bookIsbn := r.FormValue("isbn")
bookTitle := r.FormValue("title")
bookFirst:= r.FormValue("afirstname")
bookLast := r.FormValue("alastname")

var response = JsonResponse{}

if bookId == "" || bookIsbn == "" || bookTitle == "" || bookFirst == "" || bookLast == "" {
	response = JsonResponse{Type: "error", Message: "You are missing some parameters."}
} else {
	db := setupDB()

	printMessage("Inserting book into Database")

	fmt.Println("Inserting new book with ID: " + bookId + " and isbn: " + bookIsbn + " book Title: " + bookTitle + " book author " + bookFirst + bookLast)

	var lastInsertID int
	err := db.QueryRow("INSERT INTO books(id, isbn , title , afirstname , alastname) VALUES($1, $2 , $3, $4 , $5) returning id;", bookId, bookIsbn, bookTitle, bookFirst, bookLast).Scan(&lastInsertID)

// check errors
checkErr(err)

response = JsonResponse{Type: "success", Message: "The book has been inserted successfully!"}
}

json.NewEncoder(w).Encode(response)
 }


//Update a Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	var response = JsonResponse{}
	params := mux.Vars(r)
    booksID := params["id"]
    bookId := r.FormValue("id")
    bookIsbn := r.FormValue("isbn")
    bookTitle := r.FormValue("title")
    bookFirst:= r.FormValue("afirstname")
    bookLast := r.FormValue("alastname")
    if booksID == ""{
        response = JsonResponse{Type: "failed",Message: "Empty"}
        json.NewEncoder(w).Encode(response)
    }else{
	db := setupDB()
  sqlStatement := `UPDATE  books SET isbn =$2,title =$3,afirstname =$4,alastname =$5 WHERE id = $1`
  res, err := db.Query(sqlStatement, bookId,bookIsbn,bookTitle,bookFirst,bookLast)
  checkErr(err)
  fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
  fmt.Println(res)
  response = JsonResponse{Type: "success", Message: "The book  has been updated successfully"}
  json.NewEncoder(w).Encode(response)
}
}

//Delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

    bookId := params["id"]

    var response = JsonResponse{}

    if bookId == "" {
        response = JsonResponse{Type: "error", Message: "No bookID specified."}
    } else {
        db := setupDB()

        printMessage("Deleting book from DB")

        _, err := db.Exec("DELETE FROM books where id = $1", bookId)

        // check errors
        checkErr(err)

        response = JsonResponse{Type: "success", Message: "The book has been deleted successfully!"}
    }

    json.NewEncoder(w).Encode(response)
}

func deleteBooks(w http.ResponseWriter, r *http.Request){
	db := setupDB()

    printMessage("Deleting all books...")

    _, err := db.Exec("DELETE FROM books database")

    // check errors
    checkErr(err)

    printMessage("All books have been deleted successfully!")

    var response = JsonResponse{Type: "success", Message: "All books have been deleted successfully!"}

    json.NewEncoder(w).Encode(response)
}

func main() {
	//Init the Mux Router
	r := mux.NewRouter()
	//Route Handlers / Endpoints
	 r.HandleFunc("/api/books", getBooks).Methods("GET")
	 r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	 r.HandleFunc("/api/books/insert", createBook).Methods("POST")
	 r.HandleFunc("/api/books/update/{id}", updateBook).Methods("PUT")
	 r.HandleFunc("/api/books/delete/  {id}", deleteBook).Methods("DELETE")
	 r.HandleFunc("/api/books", deleteBooks).Methods("DELETE")
	fmt.Println("Server Started Successfully!")
	log.Fatal(http.ListenAndServe(":8000", r))
}
