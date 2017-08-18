package main

//Author: C Neuhardt
//Last Updated: 8/17/2017
//Last Updated By: Ricky Clevinger
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
