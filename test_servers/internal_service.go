package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`)
	})
	http.ListenAndServe(":8081", nil)
}
