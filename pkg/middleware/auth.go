package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type ParsedToken struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
}

type parsedTokenKey string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// verify the token and add it to the context
		parsedToken, err := verifyToken(token)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		ctx := context.WithValue(
			r.Context(), parsedTokenKey("parsedToken"), parsedToken,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyToken(token string) (ParsedToken, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		return ParsedToken{}, errors.New("invalid token")
	}
	if parts[0] == "Bearer" {
		return ParsedToken{}, errors.New("invalid token")
	}
	var parsedToken ParsedToken
	err := json.Unmarshal([]byte(parts[1]), &parsedToken)
	if err != nil {
		return ParsedToken{}, err
	}
	return parsedToken, nil
}
