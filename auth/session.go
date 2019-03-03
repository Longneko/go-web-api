package auth

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "net/http"
    "os"
    "../files"
)

const (
    IdByteLen = 16
    SessionStoragePath = "./_temp_sessions"
    SessionIdCookieName = "session_id"
    SessionIdCookieMaxAge = 86400
)

type session struct{
    id, username string
}

// generateRandomId generates a random 16 bit id and returns it as a hex encoded string.
// Currently used for generating session ids
func generateRandomId() (string, error) {
    b := make([]byte, IdByteLen)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b[:]), nil
}

// NewSession accepts a *user pointer and returns a pointer to a new session for that user 
func NewSession(u *user) (*session, error) {
    id, err := generateRandomId()
    if err != nil {
        return nil, err
    }
    return &session{id, u.GetUsername()}, nil
}

// sessionFromFile accepts an id string and returns a session of said id pointer constructed from 
// the file. See *session.Write() for file related details
func sessionFromFile(id string) (*session, error) {
    filepath := SessionStoragePath + "/" + id
    lines, err := files.ScanFileByLines(filepath)
    if err != nil {
        return nil, err
    }

    return &session{id, lines[0]}, nil
}

// Write creates a file named after the session id in the SessionStoragePath directory. Session 
// values are stored in separate lines in order of their appearance in the stuct (skipping the id
// as it is used for the file name)
func (s *session) Write() error {
    filepath := SessionStoragePath + "/" + s.id
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer file.Close()

    if _, err := file.WriteString(s.username + "\n"); err != nil {
        return err
    }

    return nil
}

// CreateCookie returns a http.Cookie pointer for tracking user session.
// Cookie's name and max age are set by the SessionIdCookieName and SessionIdCookieMaxAge constants
// respectively
func (s *session) CreateCookie() *http.Cookie {
    id := s.GetId()
    raw := SessionIdCookieName + "=" + id
    sessionCookie := http.Cookie{
        Name    : SessionIdCookieName,
        Value   : id,
        MaxAge  : SessionIdCookieMaxAge,
        Secure  : true,
        HttpOnly: true,
        Raw     : raw,
        Unparsed: []string{raw},
    }

    return &sessionCookie
}

// SessionFromCookie accepts a cookie with session id and returns constructed session if found.
// Cookie name must match SessionIdCookieName constant
func SessionFromCookie(cookie *http.Cookie) (*session, error) {
    if cookie.Name != SessionIdCookieName {
        err := fmt.Errorf("Invalid cookie name '%s'. Must be %s", cookie.Name, SessionIdCookieName)
        return nil, err
    }
    sessionId := cookie.Value
    return sessionFromFile(sessionId)
}

func (s *session) GetId() string {
    return s.id
}
