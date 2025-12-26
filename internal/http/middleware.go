package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"secure-service/internal/auth"
	"strings"
)

type contextKey string

const bearerPrefix string = "Bearer "
const maxBodySize int = 1 << 20

const parsedBodyKey contextKey = "parsedBody"
const userIDKey contextKey = "userID"

func ApiCommonMiddleware(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			log.Printf("Method %s not allowed", r.Method)
			sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next(w, r)
	}
}

func ModelBodyMiddleware[T any](next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))
			defer r.Body.Close()

			var data T
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&data); err != nil {
				log.Println(err)
				if errors.As(err, new(*http.MaxBytesError)) {
					sendErrorResponse(w, "Request body too large", http.StatusRequestEntityTooLarge)
					return
				}
				sendErrorResponse(w, "Could not parse body as json", http.StatusBadRequest)
				return
			}
			if decoder.More() {
				sendErrorResponse(w, "More than one json object in payload", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), parsedBodyKey, &data)
			r = r.WithContext(ctx)
		}

		next(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendErrorResponse(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			sendErrorResponse(w, "Invalid authorization scheme", http.StatusUnauthorized)
			return
		}

		reqToken := strings.TrimPrefix(authHeader, bearerPrefix)
		if reqToken == "" {
			sendErrorResponse(w, "Auth token not provided", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(reqToken)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Auth token invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int(claims.UserID))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
