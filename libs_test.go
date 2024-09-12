package main

import (
	"libcrypto"
	"libzk13"
	"testing"

	"github.com/kortschak/utter"
)

func TestNetworkAddressZKP(t *testing.T) {
	// Test coordinates
	testCases := []struct {
		lat float64
		lon float64
	}{
		{37.8199, -122.4783},
		{40.7128, -74.0060}, // New York
		{51.5074, -0.1278},  // London
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			// Generate address
			address, err := libcrypto.GenerateAddress(tc.lat, tc.lon, 256)
			if err != nil {
				t.Fatalf("Failed to generate address: %v", err)
			}
			s := utter.Sdump(address)
			t.Log(s)
			_, err = libzk13.NewZK13(address.ZKPProof, 256)
			if err != nil {
				t.Fatalf("failed to generate new zk13: %v", err)
			}

			//r, p := zkp.Prover()

			/*err = libzk13.NewVerifier()
			if err != nil {
				t.Errorf("ZKP validation failed for coordinates (%f, %f)", tc.lat, tc.lon)

				t.Fatalf("Error verifying address proof: %v", err)
			}*/

		})
	}
}
