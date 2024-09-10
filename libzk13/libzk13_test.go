package libzk13

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZK13(t *testing.T) {
	secretBaggage := "test_secret"
	bits := 256 // Reduced from 2048 for faster testing. Use larger bit size in production.

	t.Run("BasicFunctionality", func(t *testing.T) {
		zk, err := NewZK13(secretBaggage, bits)
		require.NoError(t, err, "Failed to create ZK13 instance")

		r, P := zk.Prover()
		assert.True(t, zk.Verifier(r, P), "Verification failed")
	})

	t.Run("DifferentSecrets", func(t *testing.T) {
		zk1, err := NewZK13("secret1", bits)
		require.NoError(t, err, "Failed to create ZK13 instance")

		zk2, err := NewZK13("secret2", bits)
		require.NoError(t, err, "Failed to create ZK13 instance")

		r, P := zk1.Prover()
		assert.False(t, zk2.Verifier(r, P), "Verification succeeded with different secrets")
	})

	t.Run("InvalidProof", func(t *testing.T) {
		zk, err := NewZK13(secretBaggage, bits)
		require.NoError(t, err, "Failed to create ZK13 instance")

		r, P := zk.Prover()
		P.Add(P, P) // Modify the proof
		assert.False(t, zk.Verifier(r, P), "Verification succeeded with invalid proof")
	})
}

func BenchmarkZK13(b *testing.B) {
	secretBaggage := "benchmark_secret"
	bits := 256 // Use smaller bit size for benchmarking

	zk, err := NewZK13(secretBaggage, bits)
	if err != nil {
		b.Fatalf("Failed to create ZK13 instance: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, P := zk.Prover()
		if !zk.Verifier(r, P) {
			b.Fatalf("Prover-Verifier interaction failed")
		}
	}
}

func BenchmarkZK13WithDifferentBitSizes(b *testing.B) {
	secretBaggage := "benchmark_secret"
	bitSizes := []int{256, 512, 1024, 1444 /* 2044, 2049*/}

	for _, bits := range bitSizes {
		b.Run(fmt.Sprintf("Bits-%d", bits), func(b *testing.B) {
			zk, err := NewZK13(secretBaggage, bits)
			require.NoError(b, err, "Failed to create ZK13 instance")

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				r, P := zk.Prover()
				if !zk.Verifier(r, P) {
					b.Fatalf("Prover-Verifier interaction failed")
				}
			}
		})
	}
}
