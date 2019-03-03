package files

import (
    "bufio"
    "os"
)

const (
    ReadAndWriteMode = 0644
    ReadOnlyMode = 0444
)

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
        return nil, err
    }

    return lines, nil
}