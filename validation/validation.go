package validation

import (
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"errors"
    "owlhnode/database"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    logs.Info("NEW HASH PASSWD--> "+string(bytes))
    return string(bytes), err
}

func CheckPasswordHash(password string, hash string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))	
	if err != nil {logs.Error(err); return false, err}
    return true, nil
}

// Generates a jwt token.
func Encode(secret string) (val string, err error) {
	claims := jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "OwlH",
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {logs.Error("Encode error: %s", err); return "", err}
	return tokenString, err
}

func CheckToken(token string)(err error){
	users,err := ndb.GetLoginData()
	for x := range users{
		tkn, err := Encode(users[x]["secret"])
		if err != nil {
			logs.Error("Error checking Master token: %s", err); return err
		}else{
			if token == tkn {
				return nil
			}else{
				return errors.New("The token retrieved is false")
			}
		}		
	}
	return errors.New("There are not token. Error checking token Token")
}