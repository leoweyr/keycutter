package entropy

import (
	"crypto/rand"

	"go.leoweyr.com/tokenforge/go/internal/encoding"
)

// Generator produces a fresh high-intensity entropy segment on demand.
type Generator interface {
	Generate() (*Entropy, error)
}

// SecureGenerator draws entropy from an operating-system CSPRNG and applies
// rejection sampling to guarantee an unbiased uniform symbol distribution.
type SecureGenerator struct {
	base62Alphabet *encoding.Alphabet
}

// NewSecureGenerator builds a SecureGenerator bound to the given Base62 alphabet.
func NewSecureGenerator(base62Alphabet *encoding.Alphabet) *SecureGenerator {
	return &SecureGenerator{base62Alphabet: base62Alphabet}
}

// drawIndex returns a single alphabet index sampled uniformly across the alphabet,
// rejecting any raw byte that would otherwise introduce modulo bias.
func (secureGenerator *SecureGenerator) drawIndex() (int, error) {
	var size int = secureGenerator.base62Alphabet.Size()
	var ceiling int = 256 - (256 % size)
	var buffer []byte = make([]byte, 1)

	for {
		var readError error
		_, readError = rand.Read(buffer)

		if readError != nil {
			return 0, readError
		}

		var candidate int = int(buffer[0])

		if candidate < ceiling {
			return candidate % size, nil
		}
	}
}

// Generate produces a fixed-width entropy segment by drawing one unbiased index per
// character position from the bound alphabet.
func (secureGenerator *SecureGenerator) Generate() (*Entropy, error) {
	var characters []byte = make([]byte, Length)

	var position int

	for position = 0; position < Length; position++ {
		var index int
		var drawError error
		index, drawError = secureGenerator.drawIndex()

		if drawError != nil {
			return nil, drawError
		}

		characters[position] = secureGenerator.base62Alphabet.CharacterAt(index)
	}

	return NewEntropy(string(characters)), nil
}
