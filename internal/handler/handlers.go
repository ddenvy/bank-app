package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"bank-app/internal/service"
)

type Handler struct {
	services *service.Services
	logger   *logrus.Logger
}

func NewHandlers(services *service.Services, logger *logrus.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *Handler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(response{Success: code < 400, Data: data})
	}
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, response{Success: false, Error: err.Error()})
}

// AuthMiddleware проверяет JWT токен
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.error(w, r, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your-secret-key"), nil // В реальном приложении брать из конфига
		})

		if err != nil {
			h.error(w, r, http.StatusUnauthorized, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
			if err != nil {
				h.error(w, r, http.StatusUnauthorized, errors.New("invalid user id"))
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			h.error(w, r, http.StatusUnauthorized, errors.New("invalid token"))
		}
	})
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register обработчик регистрации пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Users.Register(r.Context(), req.Username, req.Email, req.Password); err != nil {
		h.error(w, r, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, r, http.StatusCreated, nil)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login обработчик аутентификации
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.error(w, r, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.Users.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.error(w, r, http.StatusUnauthorized, err)
		return
	}

	h.respond(w, r, http.StatusOK, loginResponse{Token: token})
}

// CreateAccount обработчик создания счета
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	if err := h.services.Accounts.Create(r.Context(), userID); err != nil {
		h.error(w, r, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, r, http.StatusCreated, nil)
}

// GetAccounts обработчик получения списка счетов
func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	accounts, err := h.services.Accounts.GetByUserID(r.Context(), userID)
	if err != nil {
		h.error(w, r, http.StatusInternalServerError, err)
		return
	}

	h.respond(w, r, http.StatusOK, accounts)
}

// GetAccount обработчик получения информации о счете
func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.error(w, r, http.StatusBadRequest, errors.New("invalid account id"))
		return
	}

	account, err := h.services.Accounts.GetByID(r.Context(), accountID)
	if err != nil {
		h.error(w, r, http.StatusNotFound, err)
		return
	}

	h.respond(w, r, http.StatusOK, account)
}

// CreateCard обработчик создания карты
func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Create card handler")
}

// GetCards обработчик получения списка карт
func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get cards handler")
}

// GetCard обработчик получения информации о карте
func (h *Handler) GetCard(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get card handler")
}

// CreateTransfer обработчик создания перевода
func (h *Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Create transfer handler")
}

// CreateCredit обработчик создания кредита
func (h *Handler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Create credit handler")
}

// GetCreditSchedule обработчик получения графика платежей
func (h *Handler) GetCreditSchedule(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get credit schedule handler")
}

// GetAnalytics обработчик получения аналитики
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Get analytics handler")
}

// PredictBalance обработчик прогноза баланса
func (h *Handler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Predict balance handler")
}
