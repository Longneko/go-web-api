package auth

import (
    "crypto/sha256"
    "encoding/csv"
    "encoding/hex"
    "errors"
    "fmt"
    "os"
    "unicode/utf8"
    "../files"
)

const (
    readAndWriteMode = 0644
    readOnlyMode = 0444
    UsersFile = "users.csv"
    UsernameLenMin = 8
    PasswordLenMin = 8
    UsernameLenMax = 40
    PasswordLenMax = 40
)

// generatePasswordHash accepts a password string and returns a hex encoded sha256 hash as a string
func generatePasswordHash(password string) string {
    sum := sha256.Sum256([]byte(password))
    return hex.EncodeToString(sum[:])
}

type user struct {
    username, passwordHash, FirstName, LastName string
}

// NewUser accepts username and password strings, returns pointer to a newly constructed user.
// FirstName and LastName are set to their zero string values and should be set outside of this 
// constructor.
func NewUser(username, password string) (*user, error) {
    u := &user{}
    if err := u.SetUsername(username); err != nil {
        return nil, err
    }
    if err := u.SetPassword(password); err != nil  {
        return nil, err
    }
    return u, nil
}

// FetchUser accepts a username string and searches through the UsersFile for the correspoinding user.
// Returns pointer to user and error. If an error occurs, returns nill and the error.
// Both values are returned as nil if EOF reached without errors
func FetchUser(username string) (*user, error) {
    file, err := os.OpenFile(UsersFile, os.O_RDONLY, files.ReadOnlyMode)
    if err != nil {
        return nil, err
    }
    reader := csv.NewReader(file)

    rows, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    for _, row := range(rows) {
        if row[0] == username {
            return &user{row[0], row[1], row[2], row[3]}, nil
        }
    }

    return nil, nil
}

// Write creates a record based on user's attribute values and appends it to the UsersFile.
// The user must have a unique username that does not exist in the UsersFile.
func (u *user) Write() error {
    exists, err := FetchUser(u.username)
    if exists != nil {
        errMsg := fmt.Sprintf("User with username \"%s\" already exists", u.username)
        return errors.New(errMsg)
    } else if err != nil {
        return err
    }

    file, err := os.OpenFile(UsersFile, os.O_APPEND|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer file.Close()
    writer := csv.NewWriter(file)
    
    record := []string{u.username, u.passwordHash, u.FirstName, u.LastName}
    if err := writer.Write(record); err != nil {
        return err
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        return err
    }

    return nil
}
// SetUsername accepts a username string and sets the correpsoding field of the user. Must be at
// least UsernameLenMin characters long
func (u *user) SetUsername(username string) error {
    length := utf8.RuneCountInString(username)
    if length < UsernameLenMin || length > UsernameLenMax {
        errMsg := fmt.Sprintf("Username must be len characters long, where %d<=len<=%d",
                              UsernameLenMin, UsernameLenMax)
        return errors.New(errMsg)
    }
    u.username = username

    return nil
}

func (u *user) GetUsername() string {
    return u.username
}

// CheckPassword accepts password string and checks its hash against the user's hash stored in
// UsersFile. Returns true if hashes match
func (u *user) CheckPassword(password string) bool {
    return u.passwordHash == generatePasswordHash(password)
}

// SetPassword accepts a password string and sets its sha256 sum  as  user's 'passwordHash'
// attribute. Password has to be at least PasswordLenMin characters long
func (u *user) SetPassword(password string) error {
    length := utf8.RuneCountInString(password)
    if length < PasswordLenMin || length > PasswordLenMax {
        errMsg := fmt.Sprintf("Password must be len characters long, where %d<=len<=%d",
                              PasswordLenMin, PasswordLenMax)
        return errors.New(errMsg)
    }
    u.passwordHash = generatePasswordHash(password)

    return nil
}


// InitUsers creates a new .csv file (named by the UsersFile constant) if able. Sets the first
// row values equal to the User struct attribute names in the corresponding order
func InitUsers() error {
    file, err := os.OpenFile(UsersFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    headers := []string{"username", "passwordHash", "FirstName", "LastName"}
    if err := writer.Write(headers); err != nil {
        errMsg := fmt.Sprintf("error writing record to csv: %s", err)
        return errors.New(errMsg)
    }
    writer.Flush()

    return nil
}
