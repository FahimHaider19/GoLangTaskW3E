package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	// "sync"
)

type StudentManager struct {
	studentData map[string]Student
	// mu          sync.Mutex
}

func NewStudentManager() *StudentManager {
	return &StudentManager{
		studentData: make(map[string]Student),
	}
}

func (sm *StudentManager) IndexRequestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (sm *StudentManager) AddStudentRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // file size limit 10mb
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if _, exists := sm.studentData[id]; exists {
		http.Error(w, "Student with this ID already exists", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	cgpa := r.FormValue("cgpa")
	career := r.FormValue("career")

	cgpaFloat := 0.0
	fmt.Sscanf(cgpa, "%f", &cgpaFloat)

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error reading image file", http.StatusBadRequest)
		return
	}

	// Create a new channel to communicate the result
	resultChan := make(chan error)

	// Prepare the data struct
	data := AddStudent{
		ID:             id,
		Name:           name,
		CGPA:           cgpaFloat,
		CareerInterest: career,
		File:           file,
		W:              w,
		ResultChan:     resultChan, // Pass the channel to the data handler
	}

	// Start the data handler in a goroutine
	go sm.AddStudentDataHandler(data)

	// Wait for the result from the data handler
	err = <-resultChan
	if err != nil {
		http.Error(w, "Error adding student data", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

func (sm *StudentManager) AddStudentDataHandler(data AddStudent) {
	imagePath := "images/" + data.ID + ".jpg"
	out, err := os.Create(imagePath)
	if err != nil {
		data.ResultChan <- err // Send the error to the result channel
		return
	}

	_, err = io.Copy(out, data.File)
	if err != nil {
		data.ResultChan <- err // Send the error to the result channel
		return
	}

	imageURL := "/images/" + data.ID + ".jpg"
	data.ImageURL = imageURL

	sm.studentData[data.ID] = Student{
		ID:             data.ID,
		Name:           data.Name,
		CGPA:           data.CGPA,
		CareerInterest: data.CareerInterest,
		ImageURL:       data.ImageURL,
	}

	data.ResultChan <- nil // Send nil to indicate success
}

func (sm *StudentManager) GetStudentRequestHandler(w http.ResponseWriter, r *http.Request) {
	// sm.mu.Lock()
	// defer sm.mu.Unlock()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	student, exists := sm.studentData[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("template/card.html")
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, student)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (sm *StudentManager) DeleteStudentRequestHandler(w http.ResponseWriter, r *http.Request) {
	// sm.mu.Lock()
	// defer sm.mu.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	_, exists := sm.studentData[id]
	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Delete the image file
	imagePath := "images/" + id + ".jpg"
	err = os.Remove(imagePath)
	if err != nil {
		http.Error(w, "Error deleting image", http.StatusInternalServerError)
		return
	}

	delete(sm.studentData, id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
