package signatures

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Storage interface {
	CreateDevice(ctx context.Context, input CreateDeviceInput) (Device, error)
	FindDevice(ctx context.Context, key uuid.UUID) (Device, error)
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error)
}

type Handler struct {
	mux     *http.ServeMux
	storage Storage
}

func NewHandler(m *http.ServeMux, s Storage) *Handler {
	handler := &Handler{mux: m, storage: s}

	handler.mux.HandleFunc("POST /device", handler.CreateDevice)
	handler.mux.HandleFunc("POST /transaction", handler.SignTransaction)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

type CreateDeviceRequest struct {
	Key       uuid.UUID `json:"key"`
	Algorithm Algorithm `json:"algorithm"`
	Label     string    `json:"label"`
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, ErrEmptyJSONBody.Error(), http.StatusBadRequest)
		return
	}

	var body CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := CreateDeviceInput{
		Key:       body.Key,
		Algorithm: body.Algorithm,
		Label:     body.Label,
	}

	// TODO: handling multiple cases
	if _, err := h.storage.CreateDevice(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

type SignTransactionRequest struct {
	DeviceID uuid.UUID `json:"deviceID"`
	Data     string    `json:"data"`
}

type SignTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signedData"`
}

func (h *Handler) SignTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, ErrEmptyJSONBody.Error(), http.StatusBadRequest)
		return
	}

	var body SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := CreateTransactionInput{
		DeviceID: body.DeviceID,
		Data:     body.Data,
	}

	transaction, err := h.storage.CreateTransaction(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(transaction)
	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
