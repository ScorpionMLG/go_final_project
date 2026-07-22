package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PasswordHash string
	jwt.RegisteredClaims
}

type Password struct {
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func CreateToken(password string) (string, error) {
	claims := Claims{PasswordHash: hashPassword(password)}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := newToken.SignedString([]byte(password))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func authHandler(res http.ResponseWriter, req *http.Request) {
	var passwordJSON Password
	err := json.NewDecoder(req.Body).Decode(&passwordJSON)
	defer req.Body.Close()
	if err != nil {
		err = fmt.Errorf("incorrect JSON. %w", err)
		writeJson(res, err)
		return
	}
	if passwordJSON.Password == os.Getenv("TODO_PASSWORD") {
		token, err := CreateToken(passwordJSON.Password)
		if err != nil {
			err = fmt.Errorf("create token error. %w", err)
			writeJson(res, err)
			return
		}
		var tokenStruct Token
		tokenStruct.Token = token
		writeJson(res, tokenStruct)
		return
	}
	writeJson(res, fmt.Errorf("Неверный пароль"))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(res, req)
			return
		}
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := req.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			var claims Claims
			token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
				// проверяем алгоритм подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неверный метод подписи")
				}
				return []byte(pass), nil
			})
			if err != nil || claims.PasswordHash != hashPassword(pass) {
				http.Error(res, "bad token.", http.StatusUnauthorized)
				return
			}
			valid = token.Valid
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(res, "authentification required.", http.StatusUnauthorized)
				return
			}
		}
		next(res, req)
	})
}
