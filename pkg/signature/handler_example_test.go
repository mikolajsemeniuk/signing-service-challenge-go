package signature_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/signature"
	"github.com/google/uuid"
)

func ExampleHandler_ListDevices() {
	// Mock storage with predefined devices
	store := &storage{
		listDevices: func(_ context.Context) ([]signature.Device, error) {
			return []signature.Device{
				{
					Key:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					Label:     "Device 1",
					Algorithm: signature.Algorithm("RSA"),
				},
				{
					Key:       uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					Label:     "Device 2",
					Algorithm: signature.Algorithm("ECC"),
				},
			}, nil
		},
	}

	// Initialize the handler with the mock storage
	handler := signature.NewHandler(store)

	// Create a new HTTP GET request to the /device endpoint
	request := httptest.NewRequest(http.MethodGet, "/device", nil)
	recorder := httptest.NewRecorder()

	// Serve the HTTP request
	handler.ServeHTTP(recorder, request)

	// Print the response body
	fmt.Println(recorder.Body.String())

	// Output:
	// [{"key":"11111111-1111-1111-1111-111111111111","publicKey":null,"privateKey":null,"algorithm":"RSA","label":"Device 1","counter":0,"transactions":null},{"key":"22222222-2222-2222-2222-222222222222","publicKey":null,"privateKey":null,"algorithm":"ECC","label":"Device 2","counter":0,"transactions":null}]
}
