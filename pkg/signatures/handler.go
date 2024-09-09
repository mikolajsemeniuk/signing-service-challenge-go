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

func NewHandler(s Storage) *Handler {
	router := http.NewServeMux()

	handler := &Handler{router: router, storage: s}

	handler.router.HandleFunc("GET /device", handler.ListDevices)
	handler.router.HandleFunc("GET /device/{key}", handler.FindDevice)
	handler.router.HandleFunc("POST /device", handler.CreateDevice)
	handler.router.HandleFunc("POST /transaction", handler.CreateTransaction)

	return handler
}

// Handler provides API compatible with HTTP and REST standards.
type Handler struct {
	router  *http.ServeMux
	storage Storage
}

// ServeHTTP is used for joining handlers to HTTP server.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
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
