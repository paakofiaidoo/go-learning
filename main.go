package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hellWorld)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
func hellWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Hello , world!")
}
