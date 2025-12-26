package http

import (
	"log"
	"net/http"

	"secure-service/internal/auth"
	"secure-service/internal/database"
	"secure-service/internal/service"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(parsedBodyKey).(*service.RegisterRequest) // не может зафейлить после мидлвары
	if err := data.Validate(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	emailBusy, err := database.UserExistsByEmail(data.Email)
	if err != nil {
		sendErrorResponse(w, "Failed to check if email is busy", http.StatusInternalServerError)
		return
	}
	if emailBusy {
		log.Printf("User with email %s already exists", data.Email)
		sendErrorResponse(w, "Email is busy", http.StatusConflict)
		return
	}

	passwordHash, err := auth.HashPassword(data.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to hash the password", http.StatusInternalServerError)
		return
	}

	user, err := database.CreateUser(data.Email, data.Username, passwordHash)
	if err != nil {
		sendErrorResponse(w, "Failed to create a user", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, user, http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(parsedBodyKey).(*service.LoginRequest)
	if err := data.Validate(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	user, err := database.GetUserByEmail(data.Email)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch a user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Printf("User with email %s does not exist", data.Email)
		sendErrorResponse(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	if ok := auth.CheckPassword(data.Password, user.PasswordHash); !ok {
		log.Println("Password does not match")
		sendErrorResponse(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(*user)
	if err != nil {
		sendErrorResponse(w, "Failed to generate a token", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(
		w,
		service.AuthResponse{
			Token: token,
			User:  *user,
		},
		http.StatusOK,
	)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		sendErrorResponse(w, "Failed to get userID from context", http.StatusInternalServerError)
		return
	}
	user, err := database.GetUserByID(userID)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch a user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		sendErrorResponse(w, "User does not exist", http.StatusNotFound)
		return
	}
	sendJSONResponse(w, user, http.StatusOK)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if database.Db != nil {
		if err := database.Db.Ping(); err != nil {
			sendErrorResponse(w, "Database connection failed", http.StatusServiceUnavailable)
			return
		}
	}
	sendJSONResponse(
		w,
		map[string]string{
			"status":  "ok",
			"message": "Service is running",
		},
		http.StatusOK,
	)
}
