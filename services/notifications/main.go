package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Notifications service running on port 8083...")
	http.ListenAndServe(":8083", nil)
}
