package checksum

// Calculator maps a token base string to its fixed-width Base62 tail checksum.
type Calculator interface {
	Calculate(baseString string) *Checksum
}
