package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Denisowiec/FoleyBookkeeper/internal/auth"
	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/google/uuid"
)

func validateEmail(inputEmail string) bool {
	// This function checks whether the provided email looks like a proper email address
	// It might be replaced later with something more robust

	// First we count @ characters in an email address. There should be exactly one.
	if strings.Count(inputEmail, "@") != 1 {
		return false
	}

	// We then check wether there's something in front and after the @
	splitEmail := strings.Split(inputEmail, "@")
	if len(splitEmail) != 2 {
		return false
	}
	if len(splitEmail[0]) == 0 || len(splitEmail[1]) == 0 {
		return false
	}
	return true
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Handler function for user creation requests
	// The user provides a name, email address and password
	// A secure connection to the frontend is asssumed, so the passowrd is hashed at the back-end
	type userInputType struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userInput := userInputType{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&userInput)
	if err != nil {
		respondWithError(w, "Error decoding user request", http.StatusInternalServerError, err)
		return
	}

	// Before the user is created, the email given is validated and
	// a password hash is created
	if !validateEmail(userInput.Email) {
		respondWithError(w, "Invalid email address", http.StatusBadRequest, nil)
		return
	}
	hashedPassword, err := auth.HashPassword(userInput.Password)
	if err != nil {
		respondWithError(w, "Error hashing user password", http.StatusInternalServerError, err)
		return
	}

	createUserParams := db.CreateUserParams{
		Username:       userInput.Username,
		Email:          userInput.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, "Error creating user", http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusCreated)

	// We respond with all the user's data, except the password hash
	type createUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
	}
	resp := createUserResponse{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
	}

	dat, err := json.Marshal(resp)
	// If there's an error at this stage, the user has already been created
	// We send back an empty payload
	if err != nil {
		log.Println("Error marshalling user response", err)
		w.Write([]byte{})
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type userInputType struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userInput := userInputType{}
	err := decoder.Decode(&userInput)
	if err != nil {
		respondWithError(w, "Error decoding user input", http.StatusBadRequest, err)
		return
	}
	if userInput.Email == "" {
		respondWithError(w, "Email address required", http.StatusBadRequest, nil)
		return
	}
	if userInput.Password == "" {
		respondWithError(w, "Password requires", http.StatusBadRequest, nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), userInput.Email)
	if err != nil {
		respondWithError(w, "User not found", http.StatusUnauthorized, err)
		return
	}

	// This checks if the passwords match
	match, err := auth.CheckPasswordHash(userInput.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, "Error hashing password", http.StatusBadRequest, err)
		return
	}
	if !match {
		respondWithError(w, "Password incorrect", http.StatusUnauthorized, nil)
		return
	}

	// Everything checks out so far, so we create a JWT for the user
	// Default expiration time should be in the apiConfig
	jwt, err := auth.MakeJWT(user.ID, cfg.secret, cfg.jwtExpirationTime)
	if err != nil {
		respondWithError(w, "Error creating a JWT", http.StatusUnauthorized, err)
		return
	}

	// We also create refresh tokens for long-term login sessions
	// Default expiration time should be in the apiConfig
	refToken := auth.MakeRefreshToken()
	refTokenExpirationDate := time.Now()
	refTokenExpirationDate.Add(cfg.refTokenExpirationTime)
	setRefTokenParams := db.SetRefTokenParams{
		Token:     refToken,
		UserID:    user.ID,
		ExpiresAt: refTokenExpirationDate,
	}
	_, err = cfg.db.SetRefToken(r.Context(), setRefTokenParams)
	if err != nil {
		respondWithError(w, "Error creating a refresh token", http.StatusInternalServerError, err)
		return
	}

	// User information is sent back as a response to a successful login attempt
	type loginResponseType struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		JWT          string    `json:"jwt"`
		RefreshToken string    `json:"refresh_token"`
	}
	loginResponse := loginResponseType{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Username:     user.Username,
		Email:        user.Email,
		JWT:          jwt,
		RefreshToken: refToken,
	}

	w.WriteHeader(http.StatusAccepted)

	dat, err := json.Marshal(loginResponse)
	if err != nil {
		log.Println("Error marshalling response to login attempt:", err)
		dat = []byte{}
	}
	w.Write(dat)
}
