package error

import (
	"blog-blog-backend/core"
	"encoding/json"
	"net/http"
)

// 400....

func HandleError400(w http.ResponseWriter) {
	error400 := core.Response{Status: 400, Body: "Bad Request!"}
	err := json.NewEncoder(w).Encode(error400)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleError401(w http.ResponseWriter) {
	error401 := core.Response{Status: 401, Body: "Unauthorized User!"}
	err := json.NewEncoder(w).Encode(error401)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleError403(w http.ResponseWriter) {
	error403 := core.Response{Status: 403, Body: "Invalid Token!"}
	err := json.NewEncoder(w).Encode(error403)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleError404(w http.ResponseWriter) {
	error404 := core.Response{Status: 404, Body: "No such key found!"}
	err := json.NewEncoder(w).Encode(error404)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 500...

func HandleError500(w http.ResponseWriter) {
	error500 := core.Response{Status: 500, Body: "Internal Server Error!"}
	err := json.NewEncoder(w).Encode(error500)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 200...

func HandleSuccess200(w http.ResponseWriter, blogPost core.Blog) {
	success200 := core.Response{Status: 200, Body: blogPost}
	err := json.NewEncoder(w).Encode(success200)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleCreated201(w http.ResponseWriter) {
	created201 := core.Response{Status: 201, Body: "Created!"}
	err := json.NewEncoder(w).Encode(created201)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
