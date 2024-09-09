package signatures

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Storage interface {
	ListDevices(ctx context.Context) ([]Device, error)
	FindDevice(ctx context.Context, key uuid.UUID) (Device, error)
	CreateDevice(ctx context.Context, input CreateDeviceInput) (Device, error)
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error)
}

type Handler struct {
	mux     *http.ServeMux
	storage Storage
}

func NewHandler(m *http.ServeMux, s Storage) *Handler {
	handler := &Handler{mux: m, storage: s}

	handler.mux.HandleFunc("GET /device", handler.ListDevices)
	handler.mux.HandleFunc("GET /device/{key}", handler.FindDevice)
	handler.mux.HandleFunc("POST /device", handler.CreateDevice)
	handler.mux.HandleFunc("POST /transaction", handler.CreateTransaction)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	devices, err := h.storage.ListDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(devices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (h *Handler) FindDevice(w http.ResponseWriter, r *http.Request) {
	key, err := uuid.Parse(r.URL.Query().Get("key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.storage.FindDevice(r.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
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

	if _, err := h.storage.CreateDevice(r.Context(), CreateDeviceInput(body)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

type CreateTransactionRequest struct {
	DeviceKey uuid.UUID `json:"deviceKey"`
	Data      string    `json:"data"`
}

type CreateTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signedData"`
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, ErrEmptyJSONBody.Error(), http.StatusBadRequest)
		return
	}

	var body CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.storage.CreateTransaction(r.Context(), CreateTransactionInput(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
