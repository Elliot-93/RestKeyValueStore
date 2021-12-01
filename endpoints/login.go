package endpoints

import (
	"RestKeyValueStore/authentication"
	"RestKeyValueStore/logger"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

type LoginHandler struct{}

const LoginRoute = "/login"

func (h LoginHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch req.Method {
	case http.MethodGet:
		handleLogin(resp, req)
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}

func handleLogin(resp http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		logger.Error("Error parsing basic auth")
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}
	expectedPassword, ok := authentication.UsersAndPasswords[username]
	if !ok {
		logger.Error(fmt.Sprintf("Error Unknown User: %s", username))
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}
	if password != expectedPassword {
		logger.Error(fmt.Sprintf("Password provided is incorrect: %s", username))
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &authentication.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "UserJWTService",
		},
	}
	// Declare the token with the algorithm used for signing,
	// That is we create a token from the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT token - do this by signing the token
	// using a secure private key (jwtKey) it will create
	// {header}.{payload}.{signature}
	tokenString, err := token.SignedString(authentication.JwtKey)
	if err != nil {
		// If there is an error creating JWT return an error
		logger.Error(fmt.Sprintf("Error creating the token: %v", err))
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	jwtCookie := &http.Cookie{
		Name:    authentication.CookieName,
		Value:   tokenString,
		Expires: expirationTime,
	}

	http.SetCookie(resp, jwtCookie)
	resp.Write([]byte(fmt.Sprintf("Bearer %s", tokenString)))

	logger.Info(fmt.Sprintf("%s logged in %v", username, jwtCookie))
}
