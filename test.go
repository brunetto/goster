package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/img/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    panic(http.ListenAndServe(":17901", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<html><head></head><body><h1>Welcome Home!</h1><a href=\"/img/mapelli_progenitor_remnant.png\">Show Image!</a></body></html>")
}