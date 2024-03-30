package main

import (
	"blog-blog-backend/authentication"
	"blog-blog-backend/core"
	"blog-blog-backend/error"

	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Authenticate - check request is authenticated
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set header content type to JSON
		w.Header().Set("Content-Type", "application/json")

		//retrieve API key
		userToken := r.Header.Get("Authorization")

		// check if Authorized and valid API key
		if userToken == "" {
			error.HandleError401(w)
			return
		} else if _, ok := core.Users.M[userToken]; ok == false {
			error.HandleError403(w)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// loggingMiddleWare - log request method and url path
func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// handlePostBlog - handle posting new blogs post
func handlePostBlog(w http.ResponseWriter, r *http.Request) {
	// retrieve key
	vars := mux.Vars(r)
	key := vars["key"]

	// set header content type to JSON
	w.Header().Set("Content-Type", "application/json")

	var newBlog core.Blog // retrieve body
	if err := json.NewDecoder(r.Body).Decode(&newBlog); err != nil {
		error.HandleError400(w)
		return
	}
	defer r.Body.Close()

	// post new blogs to storage
	err := core.Post(key, newBlog)
	if err != nil {
		error.HandleError500(w)
		return
	}

	// success in post
	error.HandleCreated201(w)
	return
}

// handleGetBlog - handle getting existing blogs post
func handleGetBlog(w http.ResponseWriter, r *http.Request) {
	// retrieve key
	vars := mux.Vars(r)
	key := vars["key"]

	// set the header content type to JSON
	w.Header().Set("Content-Type", "application/json")

	blogPost, err := core.Get(key)
	// check if value exists
	if errors.Is(err, core.ErrorNoSuchKey) {
		error.HandleError404(w)
		return
	}

	// check if no error during get
	if err != nil {
		error.HandleError500(w)
		return
	}

	// send success200 back
	error.HandleSuccess200(w, blogPost)
	return
}

// handleDeleteBlog - handle deleting existing blogs post
func handleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// set header content type to JSON
	w.Header().Set("Content-Type", "application/json")

	err := core.Delete(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		error.HandleError404(w)
		return
	}
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	// set the header content type to JSON
	w.Header().Set("Content-Type", "application/json")

	var credentials core.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		error.HandleError400(w)
		return
	}
	defer r.Body.Close()

	// generate token
	generatedToken, expirationDate, err := authentication.GenerateToken(credentials.Username, credentials.Password)
	if err != nil {
		error.HandleError500(w)
		return
	}

	// send response
	formattedDate := expirationDate.Format("2006-01-02") // format date to YYYY-MM-DD
	formattedTime := expirationDate.Format("15:04:05")   // format time to HH:MM:SS
	getTokenResponse := core.Response{
		Status: 201,
		Body: core.TokenResponse{
			Token:          generatedToken,
			ExpirationDate: formattedDate,
			ExpirationTime: formattedTime,
		},
	}
	err = json.NewEncoder(w).Encode(getTokenResponse)
	if err != nil {
		error.HandleError500(w)
		return
	}

	// post credentials & string token
	err = core.PostCredentials(generatedToken, credentials)
	if err != nil {
		error.HandleError500(w)
		return
	}
}

func main() {
	r := mux.NewRouter()

	// append logging middleware function to chain
	r.Use(loggingMiddleWare)

	r.HandleFunc("/get-token", handleGetToken).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(Authenticate) // append authentication middleware function to chain
	api.HandleFunc("/post/{key}", handlePostBlog).Methods("POST")
	api.HandleFunc("/get/{key}", handleGetBlog).Methods("GET")
	api.HandleFunc("/delete/{key}", handleDeleteBlog).Methods("DELETE")

	fmt.Println("Server is running on port 8080....")
	log.Fatal(http.ListenAndServe(":8080", r))
}
