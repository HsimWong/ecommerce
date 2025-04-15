package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func TestUserRegister() {
	url := "http://localhost:8080/api/v1/user/register"
	jsonBody := []byte(`{"username":"John Doe","email":"john@example.com", "password":"password?"}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func main() {
	TestUserRegister()
}
