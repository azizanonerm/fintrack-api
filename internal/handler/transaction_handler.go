package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/azizanonerm/fintrack-api/internal/middleware"
	"github.com/azizanonerm/fintrack-api/internal/service"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	service *service.TransactionService
}

type createTransactionRequest struct {
	CategoryID      string  `json:"category_id"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDText, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDText)
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

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		http.Error(w, "transaction_date must use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Create(
		r.Context(),
		userID,
		service.CreateTransactionInput{
			CategoryID:      categoryID,
			Type:            req.Type,
			Amount:          req.Amount,
			Description:     req.Description,
			TransactionDate: transactionDate,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userIDText, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDText)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	transactions, err := h.service.List(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to list transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userIDText, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDText)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	transactionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
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

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		http.Error(w, "transaction_date must use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Update(
		r.Context(),
		userID,
		transactionID,
		service.CreateTransactionInput{
			CategoryID:      categoryID,
			Type:            req.Type,
			Amount:          req.Amount,
			Description:     req.Description,
			TransactionDate: transactionDate,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userIDText, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDText)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	transactionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), userID, transactionID); err != nil {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
