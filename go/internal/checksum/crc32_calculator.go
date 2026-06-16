package checksum

import (
	"hash/crc32"

	"go.leoweyr.com/keycutter/go/v2/internal/encoding"
)

// Crc32Calculator derives the tail checksum from the CRC32-IEEE value of a base
// string encoded as Base62.
type Crc32Calculator struct {
	encoder *encoding.Base62Codec
}

// NewCrc32Calculator builds a Crc32Calculator backed by the given Base62 codec.
func NewCrc32Calculator(encoder *encoding.Base62Codec) *Crc32Calculator {
	return &Crc32Calculator{encoder: encoder}
}

// Calculate computes the CRC32-IEEE value of the base string interpreted as a UTF-8
// byte stream and encodes it as a fixed-width Base62 checksum.
func (crc32Calculator *Crc32Calculator) Calculate(baseString string) *Checksum {
	var value uint32 = crc32.ChecksumIEEE([]byte(baseString))
	var encoded string = crc32Calculator.encoder.Encode(value, Length)

	return NewChecksum(encoded)
}
