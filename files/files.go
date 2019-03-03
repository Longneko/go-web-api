package files

import (
    "bufio"
    "os"
)

// ScanFileByLines accepts a filepath string and returns a slice of strings, each representing a
// line read from the file. Errors encountered during reading will result in returning both the 
// error and lines read so far
func ScanFileByLines(filepath string) ([]string, error) {
    var lines []string

    file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    sc := bufio.NewScanner(file)
    for sc.Scan() {
        lines = append(lines, sc.Text())
    }
    if err := sc.Err(); err != nil {
        return lines, err
    }

    return lines, nil
}