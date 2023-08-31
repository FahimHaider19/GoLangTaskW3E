package main

import (
	"mime/multipart"
	"net/http"
	// "net/http"
)

type Student struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	CGPA           float64 `json:"cgpa"`
	CareerInterest string  `json:"careerInterest"`
	ImageURL       string  `json:"imageUrl"`
}

type AddStudent struct {
	ID             int                 `json:"id"`
	Name           string              `json:"name"`
	CGPA           float64             `json:"cgpa"`
	CareerInterest string              `json:"careerInterest"`
	ImageURL       string              `json:"imageUrl"`
	File           multipart.File      `json:"file"`
	ResultChan     chan error          `json:"resultChan"`
	W              http.ResponseWriter `json:"w"`
	R              *http.Request       `json:"r"`
}

type MessageData struct {
	Message string `json:"message"`
}
