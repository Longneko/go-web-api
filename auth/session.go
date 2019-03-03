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
    SessionIdCookieDelete = "deleted"
)

type session struct{
    id string
    user *user
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

// NewSession accepts a *user pointer and returns a pointer to a new session for that user. 
// InitSession func should be used isntead unless session id cookie must be set manually or skipped
func NewSession(u *user) (*session, error) {
    id, err := generateRandomId()
    if err != nil {
        return nil, err
    }
    return &session{id, u}, nil
}

// Initsession accepts a user pointer and a http.ResponseWriter. Returns a new session pointer.
// The session is also stored in the file system
func InitSession(u *user, w http.ResponseWriter) (*session, error){
    session, err := NewSession(u)
    if err != nil {
        return nil, err
    }
    if err := session.Write(); err != nil {
        return nil, err
    }

    http.SetCookie(w, session.CreateCookie())

    return session, nil
}

// Terminate accepts a http.ResponseWriter to set the session id cookie to be deleted and removes
// the session file from server storage
func (s *session) Terminate(w http.ResponseWriter) error {
    http.SetCookie(w, GetSessionDeleteCookie())
    return s.Delete()
}

// sessionFromFile accepts an id string and returns a session of said id pointer constructed from 
// the file. See *session.Write() for file related details
func sessionFromFile(id string) (*session, error) {
    filepath := SessionStoragePath + "/" + id
    lines, err := files.ScanFileByLines(filepath)
    if err != nil {
        return nil, err
    }
    user, err := FetchUser(lines[0])
    if err != nil {
        return nil, err
    } else if user == nil {
        return nil, fmt.Errorf("Session user does not exit")
    }
    return &session{id, user}, nil
}

// Write creates a file named after the session id in the SessionStoragePath directory. Session 
// values are stored in separate lines in order of their appearance in the stuct (skipping the id
// as it is used for the file name). InitSession function incorporates this function, so it should
// not be used unless session files are managed manually
func (s *session) Write() error {
    filepath := SessionStoragePath + "/" + s.id
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer file.Close()

    if _, err := file.WriteString(s.user.username + "\n"); err != nil {
        return err
    }

    return nil
}

// Delete removes the session file and returns error if any. Session's Terminate method incorporates
// this function, so it should not be used unless session id files are managed manually
func (s *session) Delete() error {
    filepath := SessionStoragePath + "/" + s.id
    return os.Remove(filepath)
}

// CreateCookie returns a http.Cookie pointer for tracking user session.
// Cookie's name and max age are set by the SessionIdCookieName and SessionIdCookieMaxAge constants
// respectively. InitSession function incorporates this function, so it should not be used unless
// session id cookies are managed manually
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

// GetSessionDeleteCookie generates a cookie that should be set for deleting user's current session
// id cookie. Session's Terminate method incorporates this function, so it should not be used unless
// session id cookies are managed manually
func GetSessionDeleteCookie() *http.Cookie {
    id := SessionIdCookieDelete
    raw := SessionIdCookieName + "=" + id
    deleteCookie := http.Cookie{
        Name    : SessionIdCookieName,
        Value   : id,
        MaxAge  : -1,
        Secure  : true,
        HttpOnly: true,
        Raw     : raw,
        Unparsed: []string{raw},
    }

    return &deleteCookie
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

// TODO add expired sessions' files cleaning functional
// TODO add session id cookie expiration refresh on logged user activity

