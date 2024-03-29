package main

import (
	"fmt"
	"net/http"
)

//* Reading a using custom errors

func main() {
	res, err := http.Get("http://goinpracticebook.com/")

	if err != nil {
		fmt.Printf("Error here %v", err)
	}

	switch {
	case 300 <= res.StatusCode && res.StatusCode < 400:
		fmt.Println("Redirect message")
	case 400 <= res.StatusCode && res.StatusCode < 500:
		fmt.Println("Client error")
	case 500 <= res.StatusCode && res.StatusCode < 600:
		fmt.Println("Server error")
	}
}
