package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Users service running on port 8081...")
	http.ListenAndServe(":8081", nil)
}
