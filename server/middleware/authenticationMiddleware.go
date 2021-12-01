package middleware

import (
	"RestKeyValueStore/authentication"
	"RestKeyValueStore/logger"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

// Authenticate ensures a user is authenticated and sets the username on the request context
func Authenticate(nextHandler http.Handler) http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {

		// Get the JWT string from the header or cookie
		var tokenString string

		authenticationHeaderValue := req.Header.Get(authentication.HeaderName)

		if authenticationHeaderValue != "" {
			tokenString = strings.TrimPrefix(authenticationHeaderValue, "Bearer ")
		} else {
			cookie, err := req.Cookie(authentication.CookieName)
			if err != nil {
				if err == http.ErrNoCookie {
					logger.Error(fmt.Sprintf("Error no token cookie set: %v", err))
					resp.WriteHeader(http.StatusUnauthorized)
					return
				}
				logger.Error(fmt.Sprintf("Error accessing token cookie: %v", err))

				resp.WriteHeader(http.StatusBadRequest)
				return
			}

			tokenString = cookie.Value
		}

		// Initialize a new instance of `Claims`
		claims := &authentication.Claims{}
		// Parse the JWT token string
		token, err := jwt.ParseWithClaims(tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return authentication.JwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				logger.Error(fmt.Sprintf("Error signature invalid %v", err))
				resp.WriteHeader(http.StatusUnauthorized)
				return
			}
			logger.Error(fmt.Sprintf("Error processing JWT token %v", err))
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			logger.Error(fmt.Sprintf("Error invalid token %v", err))
			resp.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(req.Context(), authentication.AuthUsernameCtxKey, claims.Username)
		requestWithUser := req.WithContext(ctxWithUser)

		nextHandler.ServeHTTP(resp, requestWithUser)
	}

	return http.HandlerFunc(fn)
}
