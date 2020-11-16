# easyJWT

A light-weight JWT-manager written without any external dependencies except standard packages of Go provides features for JSON Web Token creation and validation.

### Default configurations

Default configurations could be changed by redefining these variables.

```golang
Secret       = "defaultpassphrase"
TokenRefresh = time.Hour
```

### type JWT

```golang
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
