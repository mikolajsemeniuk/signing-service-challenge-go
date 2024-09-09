package signature

// memory.go implements an in-memory storage system for managing signature devices and their transactions.
// It provides concurrency-safe operations for listing, finding, and creating devices, as well as for creating transactions.
// The in-memory store is protected by a read-write mutex to ensure thread safety, and the devices are stored using their UUID as the key.

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
	"github.com/google/uuid"
)

// Memory represents an in-memory storage for devices with concurrency control.
type Memory struct {
	mu      *sync.RWMutex
	Devices map[uuid.UUID]Device
}

// NewMemory initializes and returns a new Memory instance.
func NewMemory() *Memory {
	memory := &Memory{
		mu:      &sync.RWMutex{},
		Devices: map[uuid.UUID]Device{},
	}

	return memory
}

// ListDevices retrieves a list of all devices from the memory store.
func (m *Memory) ListDevices(_ context.Context) ([]Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var i int

	devices := make([]Device, len(m.Devices))

	for _, v := range m.Devices {
		devices[i] = v
		i++
	}

	return devices, nil
}

// FindDevice finds a device in the memory store by its UUID key.
func (m *Memory) FindDevice(_ context.Context, key uuid.UUID) (Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.Devices[key]
	if !exists {
		return device, ErrDeviceNotFound
	}

	return device, nil
}

// CreateDeviceInput holds the input data for creating a new device.
type CreateDeviceInput struct {
	Key       uuid.UUID
	Algorithm Algorithm
	Label     string
}

// CreateDevice creates a new device in the memory store.
func (m *Memory) CreateDevice(ctx context.Context, input CreateDeviceInput) (Device, error) {
	if _, err := m.FindDevice(ctx, input.Key); !errors.Is(err, ErrDeviceNotFound) {
		return Device{}, ErrDeviceAlreadyExists
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	keys := cryptic.GenerateECDSAWithMarshal
	if input.Algorithm == RSA {
		keys = cryptic.GenerateRSAWithMarshal
	}

	public, private, err := keys()
	if err != nil {
		return Device{}, nil
	}

	device := Device{
		Key:          input.Key,
		Algorithm:    input.Algorithm,
		PublicKey:    public,
		PrivateKey:   private,
		Label:        input.Label,
		Transactions: []Transaction{},
	}

	m.Devices[input.Key] = device

	return device, nil
}

// CreateTransactionInput holds the input data for creating a new transaction.
type CreateTransactionInput struct {
	DeviceKey uuid.UUID
	Data      string
}

// CreateTransaction creates a new transaction associated with a device and updates the device state.
func (m *Memory) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	device, err := m.FindDevice(ctx, input.DeviceKey)
	if err != nil {
		return Transaction{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	previous := base64.StdEncoding.EncodeToString([]byte(device.Key.String()))
	if device.Counter > 0 {
		previous = device.Transactions[len(device.Transactions)-1].Signature
	}

	data := strconv.FormatInt(device.Counter, 10) + "." + input.Data + "." + previous

	sign := cryptic.UnmarshalECDSAWithSign
	if device.Algorithm == RSA {
		sign = cryptic.UnmarshalRSAWithSign
	}

	signature, err := sign([]byte(data), device.PrivateKey)
	if err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		Signature:  string(signature),
		SignedData: data,
	}

	device.Transactions = append(device.Transactions, transaction)
	device.Counter++

	m.Devices[device.Key] = device

	return transaction, nil
}
