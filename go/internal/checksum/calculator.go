package checksum

import (
	"hash/crc32"

	"go.leoweyr.com/tokenforge/go/internal/encoding"
)

// Calculator maps a token base string to its fixed-width Base62 tail checksum.
type Calculator interface {
	Calculate(baseString string) *Checksum
}

// Crc32Calculator derives the tail checksum from the CRC32-IEEE value of a base
// string encoded as Base62.
type Crc32Calculator struct {
	encoder *encoding.Base62Encoder
}

// NewCrc32Calculator builds a Crc32Calculator backed by the given Base62 encoder.
func NewCrc32Calculator(encoder *encoding.Base62Encoder) *Crc32Calculator {
	return &Crc32Calculator{encoder: encoder}
}

// Calculate computes the CRC32-IEEE value of the base string interpreted as a UTF-8
// byte stream and encodes it as a fixed-width Base62 checksum.
func (crc32Calculator *Crc32Calculator) Calculate(baseString string) *Checksum {
	var value uint32 = crc32.ChecksumIEEE([]byte(baseString))
	var encoded string = crc32Calculator.encoder.Encode(value, Length)

	return NewChecksum(encoded)
}
