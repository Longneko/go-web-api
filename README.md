# go-web-api
This is a test assignment

## How to start
This application uses only Go's standard libraries and requires no additional installations but the language itself

### Initializing storage
The app uses filesystem to store users and sessions. So before starting the server application, it is necessary to initialize the storage file and directories. Simply run the following command while in the app's root directory:
```bash
go run init.go
```

### Running the server
The server can be started by the simple command while in the app's root directory:
```bash
go run app.go
```
This should produce console output like "Listening on :8080" if everything is ok. If the storage was not initialized, you will see the corresponding error message instead.
Currently, the server only runs on 8080 port.


## Available Routes
### /signup [POST only]
Registers new user
**Expected parameters:**
* _username_: a string to be used as username. Can only contain alphanumeric charaters (a-zA-Z) and underscore (_)
* _password_: a string to be used as password.
* _passwordConfirm_: a string, must match password.
* _firstName_: [optional] a string to be used as firstName.
* _lastName_: [optional] a string to be used as lastName.

_username_ and _password_ lengths are limited to min <= len <= max where min and max are set by the UsernameLenMin, UsernameLenMax, PasswordLenMin, PasswordLenMax constants in auth package


### /signin [POST only]
Logs the user in, sets a session id cookie.
**Expected parameters:**
_username_: a string - user's username.
_password_: a string - user's password.

### /protected
Accepts users with valid sessions only. Simply tells the user their username when accessed successfully. Returns 403 Forbidden status otherwise

### /logout
Logs the user out, removes session if able.