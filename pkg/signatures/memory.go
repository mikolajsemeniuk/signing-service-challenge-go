package signatures

import (
	"context"
	"encoding/base64"
	"sync"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
	"github.com/google/uuid"
)

type Memory struct {
	mu           *sync.RWMutex
	devices      map[uuid.UUID]Device
	transactions []Transaction
}

func NewMemory() *Memory {
	memory := &Memory{
		mu:      &sync.RWMutex{},
		devices: map[uuid.UUID]Device{},
	}

	return memory
}

func (m *Memory) FindDevice(ctx context.Context, key uuid.UUID) (Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.devices[key]
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
	m.mu.Lock()
	defer m.mu.Unlock()

	keys := crypto.GenerateECDSAWithMarshal
	if input.Algorithm == RSA {
		keys = crypto.GenerateRSAWithMarshal
	}

	public, private, err := keys()
	if err != nil {
		return Device{}, nil
	}

	device := Device{
		Key:        input.Key,
		Algorithm:  input.Algorithm,
		PublicKey:  public,
		PrivateKey: private,
		Label:      input.Label,
	}

	m.devices[input.Key] = device

	return device, nil
}

type CreateTransactionInput struct {
	DeviceID uuid.UUID
	Data     string
}

func (m *Memory) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	device, err := m.FindDevice(ctx, input.DeviceID)
	if err != nil {
		return Transaction{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	lastSignature := base64.StdEncoding.EncodeToString([]byte(device.Key.String()))
	if device.Counter > 0 {
		lastSignature = device.Transactions[len(device.Transactions)-1].Signature
	}

	data := string(device.Counter) + "." + input.Data + "." + lastSignature

	sign := crypto.UnmarshalECDSAWithSign
	if device.Algorithm == RSA {
		sign = crypto.UnmarshalRSAWithSign
	}

	signature, err := sign([]byte(data), device.PrivateKey)
	if err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{
		Signature:  string(signature),
		SignedData: data,
		Created:    time.Now(),
	}

	device.Transactions = append(device.Transactions, transaction)
	device.Counter += 1

	m.devices[device.Key] = device

	return transaction, nil
}
