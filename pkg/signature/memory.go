package signature

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
	"github.com/google/uuid"
)

type Memory struct {
	mu      *sync.RWMutex
	Devices map[uuid.UUID]Device
}

func NewMemory() *Memory {
	memory := &Memory{
		mu:      &sync.RWMutex{},
		Devices: map[uuid.UUID]Device{},
	}

	return memory
}

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

func (m *Memory) FindDevice(_ context.Context, key uuid.UUID) (Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.Devices[key]
	if !exists {
		return device, ErrDeviceNotFound
	}

	return device, nil
}

type CreateDeviceInput struct {
	Key       uuid.UUID
	Algorithm Algorithm
	Label     string
}

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

type CreateTransactionInput struct {
	DeviceKey uuid.UUID
	Data      string
}

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
