package signature_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signature"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock storage for testing.
type storage struct {
	listDevices       func(ctx context.Context) ([]signature.Device, error)
	findDevice        func(ctx context.Context, key uuid.UUID) (signature.Device, error)
	createDevice      func(ctx context.Context, input signature.CreateDeviceInput) (signature.Device, error)
	createTransaction func(ctx context.Context, input signature.CreateTransactionInput) (signature.Transaction, error)
}

func (s *storage) ListDevices(ctx context.Context) ([]signature.Device, error) {
	return s.listDevices(ctx)
}

func (s *storage) FindDevice(ctx context.Context, key uuid.UUID) (signature.Device, error) {
	return s.findDevice(ctx, key)
}

func (s *storage) CreateDevice(ctx context.Context, input signature.CreateDeviceInput) (signature.Device, error) {
	return s.createDevice(ctx, input)
}

func (s *storage) CreateTransaction(ctx context.Context, input signature.CreateTransactionInput) (signature.Transaction, error) {
	return s.createTransaction(ctx, input)
}

func TestHandler_ListDevices(t *testing.T) {
	t.Parallel()

	store := &storage{
		listDevices: func(_ context.Context) ([]signature.Device, error) {
			return []signature.Device{
				{Key: uuid.New(), Label: "Device 1", Algorithm: signature.Algorithm("RSA")},
				{Key: uuid.New(), Label: "Device 2", Algorithm: signature.Algorithm("ECC")},
			}, nil
		},
	}

	handler := signature.NewHandler(store)
	request := httptest.NewRequest(http.MethodGet, "/device", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var devices []signature.Device

	err := json.NewDecoder(recorder.Body).Decode(&devices)
	require.NoError(t, err)
	assert.Len(t, devices, 2)
}

func TestHandler_FindDevice(t *testing.T) {
	t.Parallel()

	deviceID := uuid.New()

	store := &storage{
		findDevice: func(_ context.Context, key uuid.UUID) (signature.Device, error) {
			if key == deviceID {
				return signature.Device{Key: deviceID, Label: "Device Found", Algorithm: signature.Algorithm("RSA")}, nil
			}

			return signature.Device{}, signature.ErrDeviceNotFound
		},
	}

	handler := signature.NewHandler(store)
	request := httptest.NewRequest(http.MethodGet, "/device/"+deviceID.String(), nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var device signature.Device
	err := json.NewDecoder(recorder.Body).Decode(&device)
	require.NoError(t, err)
	assert.Equal(t, deviceID, device.Key)
	assert.Equal(t, signature.Label("Device Found"), device.Label)
}

func TestHandler_CreateDevice(t *testing.T) {
	t.Parallel()

	deviceID := uuid.New()
	device := map[string]any{
		"Key":       deviceID,
		"Algorithm": signature.ECC,
		"Label":     "Test Device",
	}

	store := &storage{
		createDevice: func(_ context.Context, input signature.CreateDeviceInput) (signature.Device, error) {
			return signature.Device{
				Key:       input.Key,
				Algorithm: input.Algorithm,
				Label:     input.Label,
			}, nil
		},
	}

	handler := signature.NewHandler(store)

	body, err := json.Marshal(device)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/device", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var createdDevice signature.Device
	err = json.NewDecoder(recorder.Body).Decode(&createdDevice)
	require.NoError(t, err)
	assert.Equal(t, deviceID, createdDevice.Key)
	assert.Equal(t, signature.Label("Test Device"), createdDevice.Label)
	assert.Equal(t, signature.ECC, createdDevice.Algorithm)
}

func TestHandler_CreateTransaction(t *testing.T) {
	t.Parallel()

	transactionID := uuid.New()
	transactionRequest := map[string]any{
		"DeviceKey": transactionID,
		"Data":      "Test Data",
	}

	store := &storage{
		createTransaction: func(_ context.Context, input signature.CreateTransactionInput) (signature.Transaction, error) {
			return signature.Transaction{
				Signature:  "dummy-signature",
				SignedData: string(input.Data),
			}, nil
		},
	}

	handler := signature.NewHandler(store)

	body, err := json.Marshal(transactionRequest)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var createdTransaction signature.Transaction
	err = json.NewDecoder(recorder.Body).Decode(&createdTransaction)
	require.NoError(t, err)
	assert.Equal(t, "dummy-signature", createdTransaction.Signature)
	assert.Equal(t, "Test Data", createdTransaction.SignedData)
}
