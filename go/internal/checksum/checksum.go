package checksum

// Length is the fixed number of Base62 characters in the tail checksum segment,
// wide enough to encode any 32-bit CRC value without truncation.
const Length int = 6

// Checksum is the fixed-width Base62 tail segment encoding the CRC32 integrity
// value of a token base string.
type Checksum struct {
	value string
}

// NewChecksum wraps the given fixed-width Base62 string as a Checksum value.
func NewChecksum(value string) *Checksum {
	return &Checksum{value: value}
}

// Value returns the raw checksum characters.
func (checksum *Checksum) Value() string {
	return checksum.value
}

// Equals reports whether this checksum carries the same characters as the other checksum.
func (checksum *Checksum) Equals(other *Checksum) bool {
	return checksum.value == other.value
}

// String returns the raw checksum characters.
func (checksum *Checksum) String() string {
	return checksum.value
}
