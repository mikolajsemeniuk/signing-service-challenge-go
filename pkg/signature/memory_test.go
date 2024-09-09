package signature_test

import (
	"context"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signature"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDevice(t *testing.T) {
	t.Parallel()

	store := signature.NewMemory()
	key := uuid.New()
	expected := signature.Device{Key: key, Label: "Test Device"}
	store.Devices[key] = expected

	t.Run("Find existing device", func(t *testing.T) {
		t.Parallel()

		device, err := store.FindDevice(context.Background(), key)

		require.NoError(t, err)
		assert.Equal(t, expected, device)
	})

	t.Run("Find non-existing device", func(t *testing.T) {
		t.Parallel()

		device, err := store.FindDevice(context.Background(), uuid.New())

		require.Error(t, err)
		assert.Equal(t, signature.ErrDeviceNotFound, err)
		assert.Equal(t, signature.Device{}, device)
	})
}

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	store := signature.NewMemory()
	ctx := context.Background()

	t.Run("Successful ECDSA device creation", func(t *testing.T) {
		t.Parallel()

		deviceID := uuid.New()
		input := signature.CreateDeviceInput{
			Key:       deviceID,
			Algorithm: signature.ECC,
			Label:     "Test ECDSA Device",
		}

		device, err := store.CreateDevice(ctx, input)
		storedDevice, exists := store.Devices[deviceID]

		require.NoError(t, err)
		assert.Equal(t, deviceID, device.Key)
		assert.Equal(t, signature.ECC, device.Algorithm)
		assert.Equal(t, signature.Label("Test ECDSA Device"), device.Label)
		assert.True(t, exists)
		assert.Equal(t, device, storedDevice)
	})

	t.Run("Successful RSA device creation", func(t *testing.T) {
		t.Parallel()

		deviceID := uuid.New()
		input := signature.CreateDeviceInput{
			Key:       deviceID,
			Algorithm: signature.RSA,
			Label:     "Test RSA Device",
		}

		device, err := store.CreateDevice(ctx, input)
		storedDevice, exists := store.Devices[deviceID]

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, deviceID, device.Key)
		assert.Equal(t, signature.RSA, device.Algorithm)
		assert.Equal(t, signature.Label("Test RSA Device"), device.Label)
		assert.True(t, exists)
		assert.Equal(t, device, storedDevice)
	})
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	store := signature.NewMemory()
	ctx := context.Background()
	deviceID := uuid.New()

	_, private, _ := cryptic.GenerateECDSAWithMarshal()

	device := signature.Device{
		Key:        deviceID,
		Algorithm:  signature.ECC,
		Counter:    0,
		Label:      "Test ECDSA Device",
		PrivateKey: private,
	}
	store.Devices[deviceID] = device

	t.Run("Successful transaction creation", func(t *testing.T) {
		t.Parallel()

		input := signature.CreateTransactionInput{
			DeviceKey: deviceID,
			Data:      "transaction-data",
		}
		transaction, err := store.CreateTransaction(ctx, input)

		// Assertions
		require.NoError(t, err)
		assert.Contains(t, transaction.SignedData, "0.transaction-data")
		assert.Equal(t, int64(1), store.Devices[deviceID].Counter)
		assert.Len(t, store.Devices[deviceID].Transactions, 1)
	})

	t.Run("Device not found", func(t *testing.T) {
		t.Parallel()

		nonExistentDeviceID := uuid.New()
		input := signature.CreateTransactionInput{
			DeviceKey: nonExistentDeviceID,
			Data:      "transaction-data",
		}
		transaction, err := store.CreateTransaction(ctx, input)

		require.Error(t, err)
		assert.Equal(t, signature.ErrDeviceNotFound, err)
		assert.Equal(t, signature.Transaction{}, transaction)
	})
}
