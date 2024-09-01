package main

import (
	"log"
	"net/http"
	stdactions "usr/local/go/bin/Process"
	structs "usr/local/go/bin/Structs"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("Data.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&structs.Student{}, &structs.Messages{}, &structs.Degrees{}, &structs.Posts{})
	r := mux.NewRouter()
	r.HandleFunc("/add_student", stdactions.AddStudent).Methods("POST")
	r.HandleFunc("/search_student", stdactions.SearchStudent).Methods("GET")
	r.HandleFunc("/messages", stdactions.GetMessages).Methods("GET")
	r.HandleFunc("/new_message", stdactions.AddMessage).Methods("POST")
	r.HandleFunc("/new_post", stdactions.Post).Methods("POST")
	r.HandleFunc("/all_posts", stdactions.Posts).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
