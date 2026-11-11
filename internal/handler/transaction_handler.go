package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/azizanonerm/fintrack-api/internal/middleware"
	"github.com/azizanonerm/fintrack-api/internal/service"
	"github.com/google/uuid"
)

type TransactionHandler struct{ service *service.TransactionService }

type createTransactionRequest struct {
	CategoryID      string  `json:"category_id"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		http.Error(w, "invalid category_id", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		http.Error(w, "transaction_date must use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	t, err := h.service.Create(r.Context(), parsedUserID, service.CreateTransactionInput{CategoryID: categoryID, Type: req.Type, Amount: req.Amount, Description: req.Description, TransactionDate: date})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}
	transactions, err := h.service.List(r.Context(), parsedUserID)
	if err != nil {
		http.Error(w, "failed to list transactions", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
