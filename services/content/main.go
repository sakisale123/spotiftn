package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Content service running on port 8082...")
	http.ListenAndServe(":8082", nil)
}
