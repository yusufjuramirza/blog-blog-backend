package core

import (
	"errors"
	"sync"
)

type Response struct {
	Status int
	Body   any
}

type Blog struct {
	Title    string
	Subtitle string
	Body     string
}

type Credentials struct {
	Username string
	Password string
}

type TokenResponse struct {
	Token          string
	ExpirationDate string
	ExpirationTime string
}

var Blogs = struct {
	sync.RWMutex
	M map[string]Blog
}{M: make(map[string]Blog)}

var Users = struct {
	sync.RWMutex
	M map[string]Credentials
}{M: make(map[string]Credentials)}

var ErrorNoSuchKey = errors.New("no such key")

// Post - write data to Blogs type
func Post(key string, blogObj Blog) error {
	Blogs.Lock()
	defer Blogs.Unlock()
	b := Blog{
		Title:    blogObj.Title,
		Subtitle: blogObj.Subtitle,
		Body:     blogObj.Body,
	}
	Blogs.M[key] = b

	return nil
}

// PostCredentials - write user credentials to credentials type
func PostCredentials(key string, credObj Credentials) error {
	Users.Lock()
	defer Users.Unlock()
	u := Credentials{Username: credObj.Username, Password: credObj.Password}
	Users.M[key] = u

	return nil
}

// Get - get data from Blogs type
func Get(key string) (Blog, error) {
	Blogs.RLock()
	defer Blogs.RUnlock()
	value, ok := Blogs.M[key]

	// if "key" does not exist return ErrorNoSuchKey
	if !ok {
		return Blog{}, ErrorNoSuchKey
	}

	return value, nil
}

// Delete - delete data from Blogs type
func Delete(key string) error {
	Blogs.Lock()
	defer Blogs.Unlock()
	_, ok := Blogs.M[key]
	if ok {
		delete(Blogs.M, key)
		return nil
	}

	return ErrorNoSuchKey
}
