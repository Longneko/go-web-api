package main

import (
    "fmt"
    "log"
    "net/http"
    "./auth"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Reached Index")
    // TODO
}

// signup accepts parameters via POST requests and creates a new user.
// Expected parameters:
// username: a string to be used as username.
// password: a string to be used as password.
// passwordConfirm: a string, must match password.
// firstName: [optional] a string to be used as firstName.
// lastName: [optional] a string to be used as lastName.
// username and password lengths are limited to min <= len <= max where min and max are set by 
// UsernameLenMin, UsernameLenMax, PasswordLenMin, PasswordLenMax constants in auth package
func signup(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Add("Allowed", http.MethodPost)
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get POST form values
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    username := r.PostFormValue("username")
    password := r.PostFormValue("password")
    passwordConfirm := r.PostFormValue("passwordConfirm")
    firstName := r.PostFormValue("firstName")
    lastName := r.PostFormValue("lastName")

    if password != passwordConfirm {
        http.Error(w, "Passwords do not match", http.StatusBadRequest)
        return
    }

    user, err := auth.NewUser(username, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    user.FirstName = firstName
    user.LastName = lastName
    if err := user.Write(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User created successfully"))
}

func signin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Add("Allowed", http.MethodPost)
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get POST form values
    if err := r.ParseForm(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Check if user exists
    username := r.PostFormValue("username")
    user, err := auth.GetUser(username)
    if user == nil && err == nil {
        http.Error(w, "User not found", http.StatusForbidden)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // check password
    password := r.PostFormValue("password")
    if !user.CheckPassword(password) {
        http.Error(w, "Invalid password", http.StatusForbidden)
        return
    }
    // TODO: add login logic

    fmt.Fprintf(w, "%v\n", user.GetUsername())
}

func handleRequests() {
    http.HandleFunc("/index", index)
    http.HandleFunc("/signin", signin)
    http.HandleFunc("/signup", signup)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    handleRequests()
}
