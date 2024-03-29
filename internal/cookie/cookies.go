// This package for using cookie
package cookies

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

// Adding custom claim user ID
type Claims struct {
	jwt.RegisteredClaims
	UserID uint32
}

// Create a new cookie
func NewCookie(w http.ResponseWriter, userID uint32) {

	token, err := createToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "Token",
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}

// Validation of given token
func Validation(tokenString string) bool {

	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

func createToken(userID uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// This function return ID current user
func GetUserID(tokenString string) (uint32, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

	if err != nil || !token.Valid {
		return 0, err
	}

	fmt.Println("Token os valid", claims.UserID)
	return claims.UserID, nil
}
