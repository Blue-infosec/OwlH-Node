package validation

import (
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"errors"
    "owlhnode/database"
)

var TokenValidated string
var UserValidated string
var UuidValidated string
var users map[string]map[string]string

// Encode generates a jwt.
func Encode(uuid string, user string, secret string) (val string, err error) {

	type MyCustomClaims struct {
		Uuid string `json:"uuid"`
		User string `json:"user"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		uuid,
		user,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "OwlH",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {logs.Error(err); return "", err}
	return tokenString, err
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    logs.Info("NEW HASH PASSWD--> "+string(bytes))
    return string(bytes), err
}

func CheckPasswordHash(password string, hash string, dbUsers map[string]map[string]string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))	
	if err != nil {logs.Error(err); return false, err}
	users = dbUsers
    return true, nil
}

func CheckToken(token string, user string, uuid string)(err error){
	logs.Info(token)
	users,err := ndb.GetLoginData()
	for x := range users{
		if (x == uuid) && (users[x]["user"] == user){
			tkn, err := Encode(uuid, user, users[x]["secret"])
			if err != nil {
				logs.Error("Error checking token: %s", err); return err
			}else{
				if token == tkn {
					TokenValidated = token
					UserValidated = user
					UuidValidated = uuid
					return nil
				}else{
					return errors.New("The token retrieved is false")
				}
			}
		}
	}
	return errors.New("There are not token. Error creating Token")
}