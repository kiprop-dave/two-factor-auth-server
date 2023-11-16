package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/kiprop-dave/2fa/storage"
	twfa "github.com/kiprop-dave/2fa/twoFa"
	"github.com/kiprop-dave/2fa/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	ServerAddress string
	TwoFa         twfa.TwoFaService
	Store         storage.Storage
}

func NewServer(twoFa twfa.TwoFaService, serverAddress string, store storage.Storage) *Server {
	return &Server{
		ServerAddress: serverAddress,
		TwoFa:         twoFa,
		Store:         store,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/admin/login", makeHandler(s.handleAdminLogin)).Methods("POST")
	router.HandleFunc("/admin/register", authWrapper(s.handleAdminRegister)).Methods("POST")

	router.HandleFunc("/user/register", authWrapper(s.handleUserRegister)).Methods("POST")
	router.HandleFunc("/user/rfid-check", makeHandler(s.handleUserRfidCheck)).Methods("POST")
	router.HandleFunc("/user/two-fa", makeHandler(s.handleTwoFa)).Methods("POST")

	router.HandleFunc("/users", makeHandler(s.handleUsers)).Methods("GET")
	router.HandleFunc("/attempts", makeHandler(s.handleAttempts)).Methods("GET")

	router.HandleFunc("/check-point/register", makeHandler(s.handleCheckPointRegister)).Methods("POST")

	fmt.Println("Listening on port " + s.ServerAddress)
	log.Fatal(http.ListenAndServe(s.ServerAddress, router))
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.Store.GetUsers()
	if err != nil {
		return WriteError(w, err)
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (s *Server) handleAttempts(w http.ResponseWriter, r *http.Request) error {
	attempts, err := s.Store.GetEntryAttempts()
	if err != nil {
		return WriteError(w, err)
	}
	return WriteJSON(w, http.StatusOK, attempts)
}

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) error {
	body := new(types.LoginRequest)

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
	}

	user, err := s.Store.GetUser(bson.M{"email": body.Email})
	if err != nil {
		if err == storage.ErrNotFound {
			fmt.Println("User not found")
			return WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
				Message: "Unauthorized",
			})
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
			Message: "Unauthorized",
		})
	}

	claims := types.AdminClaims{
		User: types.UserClaims{
			Email: user.Email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "x-session-id",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Value:    tokenString,
	}

	http.SetCookie(w, &cookie)

	return WriteJSON(w, http.StatusOK, types.LoginResponse{
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *Server) createUser(role string, body *types.RegistrationRequest) (*storage.User, error) {
	user := storage.User{
		ID:       primitive.NewObjectID(),
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		TagId:    body.TagId,
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(pwd)
	twfaQr, err := s.TwoFa.GenerateTwoFa(user.Email)
	if err != nil {
		return nil, err
	}
	user.TwoFaSecret = twfaQr.Secret
	user.TwoFaQrUri = twfaQr.Uri
	user.Role = role

	err = s.Store.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Server) handleAdminRegister(w http.ResponseWriter, r *http.Request) error {
	body := new(types.RegistrationRequest)

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
		return nil
	}

	user, err := s.createUser("ADMIN", body)
	if err != nil {
		if err == storage.ErrConflict {
			WriteJSON(w, http.StatusConflict, ErrorResponse{
				Message: "User already exists",
			})
			return nil
		}
		return err
	}

	return WriteJSON(w, http.StatusCreated, types.RegistrationResponse{
		TwoFaQrUri: user.TwoFaQrUri,
		ID:         user.ID.Hex(),
	})
}

func (s *Server) handleUserRegister(w http.ResponseWriter, r *http.Request) error {
	body := new(types.RegistrationRequest)

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
		return err
	}

	user, err := s.createUser("USER", body)
	if err != nil {
		if err == storage.ErrConflict {
			WriteJSON(w, http.StatusConflict, ErrorResponse{
				Message: "User already exists",
			})
			return nil
		}
		return err
	}

	return WriteJSON(w, http.StatusCreated, types.RegistrationResponse{
		TwoFaQrUri: user.TwoFaQrUri,
		ID:         user.ID.Hex(),
	})
}

func (s *Server) handleCheckPointRegister(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("name")

	if name == "" {
		return WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
	}

	point, err := s.Store.CreateCheckPoint(name)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, point)
}

func (s *Server) handleUserRfidCheck(w http.ResponseWriter, r *http.Request) error {
	body := types.RfidCheckRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
	}

	query := bson.M{"tagId": body.TagId}
	user, err := s.Store.GetUser(query)
	if err != nil {
		if err == storage.ErrNotFound {
			return WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
				Message: "Unauthorized",
			})
		}
		return err
	}
	checkpoint, err := s.Store.GetCheckPoint(bson.M{"apiKey": body.ApiKey})
	if err != nil {
		return WriteError(w, err)
	}

	entryAttempt := storage.EntryAttempt{
		ID:           primitive.NewObjectID(),
		UserId:       user.ID,
		CheckPointId: checkpoint.ID,
		TagId:        body.TagId,
		Time:         time.Now(),
		Successful:   false,
	}
	id, err := s.Store.SaveEntryAttempt(&entryAttempt)
	if err != nil {
		return err
	}

	response := types.RfidCheckResponse{
		EntryAttemptId: id,
		Role:           user.Role,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *Server) handleTwoFa(w http.ResponseWriter, r *http.Request) error {
	body := types.TwoFaRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request",
		})
	}

	//TODO:Try to use $lookup instead of two separate queries

	attempt, err := s.Store.GetEntryAttempt(body.EntryAttemptId)
	if err != nil {
		fmt.Println("Error in getting entry attempt")
		if err == storage.ErrNotFound {
			return WriteError(w, err)
		}
	}

	query := bson.M{"tagId": attempt.TagId}
	user, err := s.Store.GetUser(query)
	if err != nil {
		fmt.Println("Error in getting user")
		return WriteError(w, err)
	}
	valid := s.TwoFa.VerifyCode(user.TwoFaSecret, body.TOTP)
	response := types.TwoFaResponse{}
	if !valid {
		return WriteJSON(w, http.StatusOK, response)
	}
	response.Success = true

	if err := s.Store.SaveSuccessfulEntryAttempt(attempt.ID.Hex()); err != nil {
		return WriteError(w, err)
	}

	return WriteJSON(w, http.StatusOK, response)
}

func authWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("x-session-id")
		if err != nil {
			fmt.Println("cookie not found")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := types.AdminClaims{}
		if err = validateToken(cookie.Value, claims); err != nil {
			fmt.Println("invalid token", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err := fn(w, r); err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func validateToken(token string, claims types.AdminClaims) error {
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) error {
	if err == storage.ErrNotFound {
		return WriteJSON(w, http.StatusNotFound, ErrorResponse{
			Message: "Not found",
		})
	}
	return WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
		Message: "Internal server error",
	})
}

type ErrorResponse struct {
	Message string `json:"message"`
}
