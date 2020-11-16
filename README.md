# easyWeb

A light-weight JWT-manager written without any external dependencies except standard packages of Go. Package provides simple features for JSON Web Token creation and validation.

### Default configurations

Default configurations could be changed by redefining these variables.

```Go
Secret       = "defaultpassphrase"
TokenRefresh = time.Hour
```

### type JWT

```Go
type JWT struct {
	User struct {
		Id   int64  `json:"id"`
		Role string `json:"role"`
	} `json:"user"`
	Expires time.Time `json:"expires"`
	Token   string    `json:"token"`
}
```

### func CreateJWT(data JWT) string

CreateJWT consumes an empty JWT struct with pre-filled User.Id, User.Role & expiration time (Expires) and then returns its an encrypted version as a string.

### func ReadJWT(value string) (JWT, bool, bool)

ReadJWT decrypts a given JSON Web Token and returns two validation booleans. First is for general validation, second is a signal for token refreshment.

### How to install

```bash
go get "github.com/azakost/easyWeb"

```

### Example of usage

```Go
package main

import (
	"fmt"
	"time"

	"github.com/azakost/easyWeb"
)

func main() {
	var data easyWeb.JWT
	data.User.Id = 123
	data.User.Role = "admin"
	data.Expires = time.Now().Add(time.Hour)
	jwt := easyJWT.CreateJWT(data)
	fmt.Println("Generated JSON web token:\n---")
	fmt.Println(jwt)
	fmt.Println("---")
	readedJWT, isValid, needToRefresh := easyJWT.ReadJWT(jwt)
	fmt.Printf("%+v\n", readedJWT)
	fmt.Printf("Is Valid: %v\n", isValid)
	fmt.Printf("Need to refresh: %v\n", needToRefresh)
}

```
