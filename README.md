# RESTful API with Go: Student Data Management

This project demonstrates the implementation of a RESTful API using Go, integrating the concepts of maps, HTTP requests, goroutines, and channels. The API allows for storing and retrieving student data, including their name, ID, CGPA, career interest, and image.

## Table of Contents

- [Part 1: Data Storage via HTTP](#part-1-data-storage-via-http)
  - [1.1. Storing Student Data](#11-storing-student-data)
  - [1.2. Handling Images](#12-handling-images)
- [Part 2: Serving Data via HTTP](#part-2-serving-data-via-http)
  - [2.1. Retrieving Student Details](#21-retrieving-student-details)
  - [2.2. Deleting Student Records](#22-deleting-student-records)
- [Part 3: Go Routines and Channels](#part-3-go-routines-and-channels)

## Features
## Part 1: Data Storage via HTTP

### 1.1 Storing Student Data

We utilize an HTML form to submit student data, including fields such as student name, ID, CGPA, career interest, and image. The data is structured as JSON or stored in global variables, ensuring uniqueness based on Student ID.

HTTP methods like POST are employed to receive and store data. We leverage Golang maps and structs for structured data storage.

### 1.2 Handling Images

For image uploads, a file upload mechanism is used. The image URLs are stored in the data storage JSON or global variable.

## Part 2: Serving Data via HTTP

### 2.1 Retrieving Student Details

A web interface includes an input field for providing a student ID. Upon submission, a card with HTML formatting displays the student details.

HTTP methods like GET are utilized to serve the requested data. Bootstrap can be applied for styling HTML cards.

### 2.2 Deleting Student Records

Another input field allows users to input a student ID for deletion. Upon submission, the corresponding record is removed from storage, and a list of student details is displayed using HTML cards.

HTTP methods like DELETE are used for record removal.

## Part 3: Go Routines and Channels

To optimize network calls and enable concurrency, Go routines and channels are implemented. This enhances the efficiency of handling multiple requests simultaneously.

### Getting Started

To run the project locally, follow these steps:

1. Clone the repository: `https://github.com/FahimHaider19/GoLangTaskW3E`
2. Navigate to the project directory: `cd GoLangTaskW3E`
3. Install dependencies: `go get`
4. Build and run the project: `go run *.go`
