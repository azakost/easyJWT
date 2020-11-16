// Package for easy JSON Web Token create and read. Basic feature of this package that it is very light-weight and written without any dependencies except standard packages of golang.
package easyJWT

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Default configurations
var (
	Secret       = "writedrunkeditsober"
	TokenLife    = time.Hour
	TokenRefresh = TokenLife / 2
)

type JWT struct {
	User struct {
		Id   int64  `json:"id"`
		Role string `json:"role"`
	} `json:"user"`
	Expires time.Time `json:"expires"`
	Token   string    `json:"token"`
}

// CreateJWT consumes an empty JWT struct with pre-filled User.Id and User.Role and then returns its an encrypted version as a string.
func CreateJWT(data JWT) string {
	exp := time.Now().Add(TokenLife)
	data.Expires = exp
	tok := []byte(strconv.FormatInt(data.User.Id, 10) + exp.String())
	data.Token = encrypt(tok)
	marshaled, errorMarshal := json.Marshal(data)
	if errorMarshal != nil {
		panic(errorMarshal)
	}
	return encrypt(marshaled)
}

// ReadJWT decrypts a given JSON Web Token and returns two validation booleans. First is for general validation, second is a signal for token refreshment.
func ReadJWT(value string) (JWT, bool, bool) {
	var blank JWT
	data, errorDecrypt := decrypt(value)
	if errorDecrypt != nil {
		return blank, false, false
	}
	errorUnmarshal := json.Unmarshal(data, &blank)
	if errorUnmarshal != nil {
		return blank, false, false
	}
	tok := strconv.FormatInt(blank.User.Id, 10) + blank.Expires.String()
	token := encrypt([]byte(tok))
	if blank.Token != token {
		fmt.Println("----- GENERATED TOKEN -------")
		fmt.Println(token)
		fmt.Println("----- TOKEN -------")
		fmt.Println(blank.Token)
		return blank, false, false
	}
	if time.Now().After(blank.Expires) {
		return blank, false, false
	}
	if blank.Expires.Before(time.Now().Add(TokenRefresh)) {
		return blank, true, true
	}
	return blank, true, false
}

func encrypt(data []byte) string {
	gcm := gcm()
	nonce := make([]byte, gcm.NonceSize())
	_, errorRead := io.ReadFull(rand.Reader, nonce)
	if errorRead != nil {
		panic(errorRead)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString([]byte(ciphertext))
}

func decrypt(value string) ([]byte, error) {
	data, errorBase64 := base64.StdEncoding.DecodeString(value)
	if errorBase64 != nil {
		return nil, errorBase64
	}
	gcm := gcm()
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	decrypted, errorOpen := gcm.Open(nil, nonce, ciphertext, nil)
	if errorOpen != nil {
		return nil, errorOpen
	}
	return decrypted, nil
}

func gcm() cipher.AEAD {
	hasher := md5.New()
	_, errorWrite := hasher.Write([]byte(Secret))
	if errorWrite != nil {
		panic(errorWrite)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	block, errorAes := aes.NewCipher([]byte(hash))
	if errorAes != nil {
		panic(errorAes)
	}
	gcm, errorGCM := cipher.NewGCM(block)
	if errorGCM != nil {
		panic(errorGCM)
	}
	return gcm
}
