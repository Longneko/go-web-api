package main

import (
    "fmt"
    "os"
    "./auth"
)

const readAndWriteMode = 0644

func main() {
    // Create session temp storage path
    if err := os.Mkdir(auth.SessionStoragePath, os.ModePerm); err != nil {
        fmt.Printf("Error creating session storage dir: '%s'\n", err.Error())
    }

    if err := auth.InitUsers(); err != nil {
        fmt.Printf("Error creating user storage file: '%s'\n", err.Error())
    }
}