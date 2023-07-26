package cookies

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

func NewCookie(w http.ResponseWriter, userID uint64) {

	token, err := createToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "Token",
		Value: token,
	}

	http.SetCookie(w, cookie)
}

func Validation(tokenString string) bool {

	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})

	if err != nil {
		return false
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return false
	}

	fmt.Println("Token os valid")
	return true
}

func createToken(userID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserId(tokenString string) (uint64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return 0, err
	}

	fmt.Println("Token os valid")
	return claims.UserID, nil
}