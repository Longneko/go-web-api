package auth

import (
    "crypto/rand"
    "encoding/hex"
    "os"
    "../files"
)

const (
    IdByteLen = 16
    SessionStoragePath = "./_temp_sessions"
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

// SessionFromFile accepts an id string and returns a session of said id pointer constructed from 
// the file. See *session.Write() for file related details
func SessionFromFile(id string) (*session, error) {
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
    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_EXCL, files.ReadAndWriteMode)
    if err != nil {
        return err
    }
    defer file.Close()

    if _, err := file.WriteString(s.username + "\n"); err != nil {
        return err
    }

    return nil
}

func (s *session) GetId() string {
    return s.id
}
