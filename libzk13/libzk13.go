package libzk13

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"

	"github.com/zeebo/blake3"
	"go.dedis.ch/kyber/v3"
)

const GENERATOR = 7

type ZK13 struct {
	p, g, Hs *big.Int
}

// NewZK13 initializes the ZK13 structure with a prime number, generator, and hashed secret.
// It addresses the correct handling of byte slices and ensures that parameters are securely generated.
func NewZK13(secretBaggage string, bits int) (*ZK13, error) {
	z := &ZK13{g: big.NewInt(GENERATOR)}
	var err error

	// Generate prime and validate parameters
	for {
		z.p, err = rand.Prime(rand.Reader, bits)
		if err != nil {
			panic(fmt.Sprintf("Failed to generate a large prime: %v", err))
		}
		if z.ValidateParameters(big.NewInt(224)) {
			break
		}
	}

	// Hash secret baggage
	hash := blake3.Sum512([]byte(secretBaggage))
	z.Hs = new(big.Int).SetBytes(hash[:])

	return z, nil
}

// calculateR calculates r = g^k mod p.
func (z *ZK13) calculateR(k *big.Int) *big.Int {
	return new(big.Int).Exp(z.g, k, z.p)
}

// calculateF calculates F = Hs*k mod (p-1).
func (z *ZK13) calculateF(k *big.Int) *big.Int {
	pMinusOne := new(big.Int).Sub(z.p, big.NewInt(1))
	return new(big.Int).Mod(new(big.Int).Mul(z.Hs, k), pMinusOne)
}

// calculateP calculates P = g^F mod p.
func (z *ZK13) calculateP(F *big.Int) *big.Int {
	return new(big.Int).Exp(z.g, F, z.p)
}

// Verify checks if the given P matches r^Hs mod p, validating the proof.
func Verify(r, Hs, P, p *big.Int) bool {
	V := new(big.Int).Exp(r, Hs, p)
	return V.Cmp(P) == 0
}

func (z *ZK13) Prover() (*big.Int, *big.Int) {
	k, _ := rand.Int(rand.Reader, z.p)
	r := z.calculateR(k)
	F := z.calculateF(k)
	P := z.calculateP(F)
	return r, P
}

func (z *ZK13) Verifier(r, P *big.Int) bool {
	V := new(big.Int).Exp(r, z.Hs, z.p)
	return subtle.ConstantTimeCompare(V.Bytes(), P.Bytes()) == 1
}

// ValidateParameters checks if the parameters meet security requirements.
func (z *ZK13) ValidateParameters(minPrimeFactor *big.Int) bool {
	pMinusOne := new(big.Int).Sub(z.p, big.NewInt(1))
	return new(big.Int).Mod(pMinusOne, z.g).Cmp(big.NewInt(0)) != 0 &&
		new(big.Int).Exp(z.g, pMinusOne, z.p).Cmp(big.NewInt(1)) == 0 &&
		new(big.Int).Div(pMinusOne, minPrimeFactor).ProbablyPrime(20)
}

// Verifier represents a zero-knowledge proof verifier
type Verifier struct {
	suite kyber.Group
}

// NewVerifier creates a new Verifier instance
func NewVerifier(suite kyber.Group) *Verifier {
	return &Verifier{
		suite: suite,
	}
}
