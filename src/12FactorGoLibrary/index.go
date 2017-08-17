package main

//Author: C Neuhardt
//Last Updated: 8/3/2017

import (
	"net/http"
	"handlers"
)


func main() {
	http.HandleFunc("/", handlers.Redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	handlers.Handles()
	http.ListenAndServe(":8080", nil)
}
