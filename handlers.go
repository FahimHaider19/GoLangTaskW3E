package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	// "sync"
)

type StudentManager struct {
	studentData map[int]Student
	// mu          sync.Mutex
}

func NewStudentManager() *StudentManager {
	return &StudentManager{
		studentData: make(map[int]Student),
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

	idForm := r.FormValue("id")
	nameForm := r.FormValue("name")
	cgpaForm := r.FormValue("cgpa")
	careerForm := r.FormValue("career")

	if idForm == "" || nameForm == "" || cgpaForm == "" || careerForm == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // file size limit 10mb
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || idForm != strconv.Itoa(id) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	cgpa, err := strconv.ParseFloat(r.FormValue("cgpa"), 64)
	if err != nil || cgpa < 0 || cgpa > 4 {
		http.Error(w, "Invalid CGPA", http.StatusBadRequest)
		return
	}

	career := r.FormValue("career")

	if _, exists := sm.studentData[id]; exists {
		http.Error(w, "Student with this ID already exists", http.StatusBadRequest)
		return
	}

	// fmt.Sscanf(cgpa, "%f", &cgpaFloat)

	file, _, err := r.FormFile("image")
	if err != nil {
		print(err)
		http.Error(w, "Error reading image file", http.StatusBadRequest)
		return
	}

	// create new channel
	resultChan := make(chan error)

	data := AddStudent{
		ID:             id,
		Name:           name,
		CGPA:           cgpa,
		CareerInterest: career,
		File:           file,
		ResultChan:     resultChan, // channel
		W:              w,
		R:              r,
	}

	// start goroutine
	go sm.AddStudentDataHandler(data)

	// get result from channel
	err = <-resultChan
	if err != nil {
		http.Error(w, "Error adding student data", http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)
	tmpl, err := template.ParseFiles("template/card.html")
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "Student added successfully")
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (sm *StudentManager) AddStudentDataHandler(data AddStudent) {
	imagePath := "images/" + strconv.Itoa(data.ID) + ".jpg"
	out, err := os.Create(imagePath)
	if err != nil {
		data.ResultChan <- err // send error to result channel
		return
	}

	_, err = io.Copy(out, data.File)
	if err != nil {
		data.ResultChan <- err // send error to result channel
		return
	}

	imageURL := "/images/" + strconv.Itoa(data.ID) + ".jpg"
	data.ImageURL = imageURL

	sm.studentData[data.ID] = Student{
		ID:             data.ID,
		Name:           data.Name,
		CGPA:           data.CGPA,
		CareerInterest: data.CareerInterest,
		ImageURL:       data.ImageURL,
	}

	data.ResultChan <- nil // nil indicate success

	http.Redirect(data.W, data.R, "/msg?Message=Student added successfully", http.StatusSeeOther)
}

func (sm *StudentManager) GetStudentRequestHandler(w http.ResponseWriter, r *http.Request) {
	// sm.mu.Lock()
	// defer sm.mu.Unlock()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idQuery := r.URL.Query().Get("id")
	if idQuery == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	studentChan := make(chan Student, 1)
	errorChan := make(chan error, 1)

	go sm.GetStudentDataHandler(id, studentChan, errorChan)

	student := <-studentChan
	err = <-errorChan

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

func (sm *StudentManager) GetStudentDataHandler(id int, studentChan chan<- Student, errorChan chan<- error) {
	student, exists := sm.studentData[id]
	if !exists {
		errorChan <- fmt.Errorf("Student not found")
		studentChan <- Student{}
		return
	}

	studentChan <- student
	errorChan <- nil
	return
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

	idForm := r.FormValue("id")
	if idForm == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idForm)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	imageDeleteChan := make(chan error, 1)
	go sm.DeleteStudentDataHandler(id, imageDeleteChan)
	err = <-imageDeleteChan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/msg?Message=Student deleted successfully", http.StatusSeeOther)
}

func (sm *StudentManager) DeleteStudentDataHandler(id int, imageDeleteChan chan<- error) {
	_, exists := sm.studentData[id]
	if !exists {
		imageDeleteChan <- fmt.Errorf("Student not found")
		return
	}

	// Delete the image file
	imagePath := "images/" + strconv.Itoa(id) + ".jpg"
	err := os.Remove(imagePath)
	if err != nil {
		imageDeleteChan <- fmt.Errorf("Error deleting image")
		return
	}

	delete(sm.studentData, id)
	imageDeleteChan <- nil
	return
}

func (sm *StudentManager) MessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("Message")

	if msg == "" {
		http.Error(w, "Error! Something went wrong", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("template/msg.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := MessageData{Message: msg}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template Message", http.StatusInternalServerError)
		return
	}
}
