package main

import (
	"mime/multipart"
	// "net/http"
)

type Student struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CGPA           float64 `json:"cgpa"`
	CareerInterest string  `json:"careerInterest"`
	ImageURL       string  `json:"imageUrl"`
}

type AddStudent struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	CGPA           float64             `json:"cgpa"`
	CareerInterest string              `json:"careerInterest"`
	ImageURL       string              `json:"imageUrl"`
	File           multipart.File      `json:"file"`
	// W              http.ResponseWriter `json:"w"`
	ResultChan     chan error          `json:"resultChan"`
}
