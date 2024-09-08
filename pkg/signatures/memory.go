package signatures

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// TODO: in-memory persistence ...
// TODO: compatible with postgres

type Memory struct {
	// TODO: Fix mutex for devices and signatures
	mu         *sync.Mutex
	devices    map[uuid.UUID]Device
	signatures map[string]struct{}
}

func NewMemory() *Memory {
	memory := &Memory{
		mu:      &sync.Mutex{},
		devices: map[uuid.UUID]Device{},
	}
	return memory
}

func (m *Memory) FindDevice(key uuid.UUID) (Device, error) {
	// TODO: Fix mutex for devices and signatures
	m.mu.Lock()
	defer m.mu.Unlock()

	device, exists := m.devices[key]
	if !exists {
		return device, ErrDeviceNotFound
	}

	return device, nil
}

type CreateDeviceInput struct {
	Key        uuid.UUID
	Algorithm  Algorithm
	PublicKey  []byte
	PrivateKey []byte
	Label      string
}

func (m *Memory) CreateDevice(ctx context.Context, input CreateDeviceInput) error {
	// TODO: Fix mutex for devices and signatures
	m.mu.Lock()
	defer m.mu.Unlock()

	key := uuid.New()
	device := Device{
		Key:        input.Key,
		Algorithm:  input.Algorithm,
		PublicKey:  input.PublicKey,
		PrivateKey: input.PrivateKey,
		Label:      input.Label,
	}

	m.devices[key] = device

	return nil
}

func (m *Memory) CreateSignature(signature string) error {
	// TODO: Fix mutex for devices and signatures
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

// type SignTransactionInput struct {
// 	Key        uuid.UUID
// 	Algorithm  Algorithm
// 	PublicKey  string
// 	PrivateKey string
// 	Label      string
// }

// func (m *Memory) SignTransaction(deviceKey uuid.UUID, data string) (Device, error) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	device, err := m.FindDevice(deviceKey)
// 	if err != nil {
// 		return device, err
// 	}

// 	return device, nil
// }
