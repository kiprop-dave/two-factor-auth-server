package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiprop-dave/2fa/storage"
	twfa "github.com/kiprop-dave/2fa/twoFa"
	"github.com/kiprop-dave/2fa/types"
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
	router.HandleFunc("/admin/register", makeHandler(s.handleAdminRegister)).Methods("POST")

	router.HandleFunc("/user/register", makeHandler(s.handleUserRegister)).Methods("POST")
	router.HandleFunc("/user/rfid-check", makeHandler(s.handleUserRfidCheck)).Methods("POST")
	router.HandleFunc("/user/two-fa", makeHandler(s.handleTwoFa)).Methods("POST")

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

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
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
		return err
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

	return WriteJSON(w, http.StatusOK, types.RegistrationResponse{
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

	return WriteJSON(w, http.StatusOK, types.RegistrationResponse{
		TwoFaQrUri: user.TwoFaQrUri,
		ID:         user.ID.Hex(),
	})
}

func (s *Server) handleCheckPointRegister(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleUserRfidCheck(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleTwoFa(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
