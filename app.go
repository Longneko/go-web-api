package main

import (
    "fmt"
    "log"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Reached Index")
    // TODO
}

func signin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Add("Allowed", http.MethodPost)
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

        return
    }
    // TODO: log user in and set cookie
}

func handleRequests() {
    http.HandleFunc("/index", index)
    http.HandleFunc("/signin", signin)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    handleRequests()
}
