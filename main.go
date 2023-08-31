package main

import (
	"fmt"
	"net/http"
)

// var studentData = make(map[string]Student)

func main() {
	deleteFilesInDirectory("images")

	sm := NewStudentManager()

	http.HandleFunc("/", sm.IndexRequestHandler)
	http.HandleFunc("/add", sm.AddStudentRequestHandler)
	http.HandleFunc("/get", sm.GetStudentRequestHandler)
	http.HandleFunc("/delete", sm.DeleteStudentRequestHandler)
	http.HandleFunc("/msg", sm.MessageHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
