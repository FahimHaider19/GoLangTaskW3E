package main

import (
	"fmt"
	"net/http"
)

// var studentData = make(map[string]Student)

func main() {
	err := deleteFilesInDirectory("images")
	fmt.Printf(`Error deleting images: %v`, err)

	sm := NewStudentManager()

	http.HandleFunc("/", sm.IndexRequestHandler)
	http.HandleFunc("/add", sm.AddStudentRequestHandler)
	http.HandleFunc("/get", sm.GetStudentRequestHandler)
	http.HandleFunc("/delete", sm.DeleteStudentRequestHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
