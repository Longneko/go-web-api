package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "./auth"
)

const (
    ListenPort = ":8080"
)

// protected accepts users with valid sessions only. Simply tells the user their username when 
// accessed successfully. Returns 403 Forbidden status otherwise
func protected(w http.ResponseWriter, r *http.Request) {    
    // Verify session
    sessionCookie, _ := r.Cookie(auth.SessionIdCookieName)
    if sessionCookie == nil {
        http.Error(w, "Session missing", http.StatusForbidden)
        return
    }
    session, _ := auth.SessionFromCookie(sessionCookie)
    if session == nil {
        http.Error(w, "Session invalid or expired", http.StatusForbidden)
        return
    }

    // Get user based on session
    user, _ := auth.UserFromSession(session)
    if user == nil {
        // session is ok, but user not found, meaning internal data is probably corrupted
        http.Error(w, "Unexpected error occurred. Please try again later",
                   http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Reached protected as: %s", user.GetUsername())))
}


// signup accepts parameters via POST requests and creates a new user.
// Expected parameters:
// username: a string to be used as username.
// password: a string to be used as password.
// passwordConfirm: a string, must match password.
// firstName: [optional] a string to be used as firstName.
// lastName: [optional] a string to be used as lastName.
// username and password lengths are limited to min <= len <= max where min and max are set by the
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
        fmt.Println(err) // Error that shouldn't be exposed to client is passed to console. Sould 
                         // eventually be replaced with propper logging
        http.Error(w, "Unexpected error occurred. Please try again later",
                   http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User created successfully"))
}

// signup accepts parameters via POST requests and logs the user in by setting a session_id cookie.
// Responds with a message of successful login or the error.
// Expected parameters:
// username: a string - user's username.
// password: a string - user's password.
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
    user, err := auth.FetchUser(username)
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
    
    // Create session and set cookie
    _, err = auth.InitSession(user, w)
    if err != nil {
        fmt.Println(err) // Error that shouldn't be exposed to client is passed to console. Sould 
                         // eventually be replaced with propper logging
        http.Error(w, "Unexpected error occurred. Please try again later",
                   http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Login successful"))
}

// logout logs the user out by setting the session_id cookie to expire and its value to 'deleted'. 
// Active session is also removed from server storage
func logout(w http.ResponseWriter, r *http.Request) {
    // Verify session
    sessionCookie, _ := r.Cookie(auth.SessionIdCookieName)
    if sessionCookie == nil {
        http.Error(w, "Session missing or already terminated", http.StatusForbidden)
        return
    }
    session, _ := auth.SessionFromCookie(sessionCookie)
    if session == nil {
        http.Error(w, "Session invalid or already terminated", http.StatusForbidden)
        return
    }
    if err := session.Terminate(w); err != nil {
        fmt.Println(err) // Error that shouldn't be exposed to client is passed to console. Sould 
                         // eventually be replaced with propper logging
        http.Error(w, "Unexpected error occurred. Please try again later",
                   http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Logout successful"))
}

func handleRequests() {
    http.HandleFunc("/protected", protected)
    http.HandleFunc("/signin", signin)
    http.HandleFunc("/signup", signup)
    http.HandleFunc("/logout", logout)
    log.Fatal(http.ListenAndServe(ListenPort, nil))
}

func main() {
    checklist := []string{auth.UsersFile, auth.SessionStoragePath}
    // Exists reports whether the named file or directory exists.
    for _, pathname := range(checklist) {
        if _, err := os.Stat(pathname); err != nil {
            if os.IsNotExist(err) {
                log.Fatal("Storage assests missing. Please run 'go run init.go' and try again")
            }
        }
    }

    fmt.Printf("Listening on %s\n", ListenPort)
    handleRequests()
}
