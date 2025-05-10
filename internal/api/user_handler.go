package api

import (
	"log"
	"github.com/melkeydev/femProject/internal/store"

)

type registerUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Bio string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *logger.Logger) *UserHandler {
	return &UserHandler {
		userStore: userStore,
		logger: logger,
	}
}

func (h *UserHandler) validateRegisterRequest(reg *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
   }
   emailRegex := regexp.MustCompiled(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
   if !emailRegex.MatchString(req.Email) {
	return errors.New("invalid email format")
   }

   if req.Password == "" {
	return errors.New("password is required")
   }

   return nil
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	
