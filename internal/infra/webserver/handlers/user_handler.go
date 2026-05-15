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
	UserDB database.UserInterface
}

type Error struct {
	Message string `json:"message"`
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpireIn := r.Context().Value("exp").(int)

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
		"exp": time.Now().Add(time.Duration(jwtExpireIn) * time.Second).Unix(),
	}

	// token := jwtauth.New("H256", []byte("secret"), nil)
	//_, tokenEncode, err := token.Encode(claims)
	//_, tokenEncode, err := h.JWT.Encode(claims)

	_, tokenEncode, err := jwt.Encode(claims)
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

// Create user godocs
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserInput true "User information"
// @Success 201
// @Failure 400
// @Failure 500 {object} Error "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	newUser, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.UserDB.Create(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
