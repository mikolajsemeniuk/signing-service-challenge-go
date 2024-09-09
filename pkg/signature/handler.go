package signature

// handler.go implements the HTTP handlers for managing signature devices and transactions.
// It provides endpoints for listing devices, finding devices by UUID, creating new devices, and creating transactions.
// These handlers interact with the underlying storage through the defined `Storage` interface, and responses in JSON format.

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// Storage defines an interface for device and transaction operations.
type Storage interface {
	ListDevices(ctx context.Context) ([]Device, error)
	FindDevice(ctx context.Context, key uuid.UUID) (Device, error)
	CreateDevice(ctx context.Context, input CreateDeviceInput) (Device, error)
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error)
}

// NewHandler creates a new HTTP handler with routing.
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

// ListDevices serves all currently stored devices.
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.storage.ListDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// FindDevice serves device with given by user key.
func (h *Handler) FindDevice(w http.ResponseWriter, r *http.Request) {
	key, err := uuid.Parse(r.PathValue("key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.storage.FindDevice(r.Context(), key)
	if errors.Is(err, ErrDeviceNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// CreateDevice saves device to datastore.
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Key       uuid.UUID `json:"key"`
		Algorithm Algorithm `json:"algorithm"`
		Label     Label     `json:"label"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.storage.CreateDevice(r.Context(), CreateDeviceInput(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// CreateTransaction saves transaction and modify device within.
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DeviceKey uuid.UUID `json:"deviceKey"`
		Data      Data      `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.storage.CreateTransaction(r.Context(), CreateTransactionInput(body))
	if errors.Is(err, ErrDeviceNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
