package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/robsonalvesdevbr/apis-go/internal/dto"
	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"github.com/robsonalvesdevbr/apis-go/internal/infra/database"
)

type UserHandler struct {
	UserDB            database.UserInterface
	JWT               *jwtauth.JWTAuth
	JWTExpireDuration int
}

func NewUserHandler(userDB database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpireDuration int) *UserHandler {
	return &UserHandler{UserDB: userDB, JWT: jwt, JWTExpireDuration: jwtExpireDuration}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var credentials dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(credentials.Email)
	if err != nil || user == nil || !user.ValidatePassword(credentials.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Duration(h.JWTExpireDuration) * time.Second).Unix(),
	}

	// token := jwtauth.New("H256", []byte("secret"), nil)

	//_, tokenEncode, err := token.Encode(claims)

	_, tokenEncode, err := h.JWT.Encode(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenEncode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
