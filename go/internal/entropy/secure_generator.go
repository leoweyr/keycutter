package entropy

import (
	"crypto/rand"

	"go.leoweyr.com/keycutter/go/v2/internal/encoding"
)

// SecureGenerator draws entropy from an operating-system CSPRNG and applies
// rejection sampling to guarantee an unbiased uniform symbol distribution.
type SecureGenerator struct {
	base62Alphabet *encoding.Alphabet
}

// NewSecureGenerator builds a SecureGenerator bound to the given Base62 alphabet.
func NewSecureGenerator(base62Alphabet *encoding.Alphabet) *SecureGenerator {
	return &SecureGenerator{base62Alphabet: base62Alphabet}
}

// Generate produces a fixed-width entropy segment by drawing unbiased indices from the bound alphabet.
// A 32-byte bulk read reduces CSPRNG syscall overhead to O(1).
func (secureGenerator *SecureGenerator) Generate() (*Entropy, error) {
	var characters []byte = make([]byte, Length)
	var size int = secureGenerator.base62Alphabet.Size()
	var ceiling int = 256 - (256 % size)
	var buffer [32]byte
	var bufferPosition int = len(buffer)

	var position int = 0

	for position < Length {
		if bufferPosition >= len(buffer) {
			var readError error
			_, readError = rand.Read(buffer[:])

			if readError != nil {
				return nil, readError
			}

			bufferPosition = 0
		}

		var candidate int = int(buffer[bufferPosition])
		bufferPosition++

		if candidate < ceiling {
			characters[position] = secureGenerator.base62Alphabet.CharacterAt(candidate % size)
			position++
		}
	}

	return NewEntropy(string(characters)), nil
}
