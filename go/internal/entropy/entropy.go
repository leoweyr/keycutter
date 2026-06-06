package entropy

// Length is the fixed number of Base62 characters in the high-intensity entropy
// segment, sized to deliver roughly 143 bits of randomness.
const Length int = 24

// Entropy is the fixed-width high-intensity random segment of a token.
type Entropy struct {
	value string
}

// NewEntropy wraps the given fixed-width random string as an Entropy value.
func NewEntropy(value string) *Entropy {
	return &Entropy{value: value}
}

// Value returns the raw entropy characters.
func (entropy *Entropy) Value() string {
	return entropy.value
}

// Length returns the number of characters in the entropy segment.
func (entropy *Entropy) Length() int {
	return len(entropy.value)
}

// String returns the raw entropy characters.
func (entropy *Entropy) String() string {
	return entropy.value
}
