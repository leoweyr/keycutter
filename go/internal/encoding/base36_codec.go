package encoding

import (
	"errors"
	"math"
)

// Base36Codec converts unsigned integers to and from variable-width Base36 strings.
// Its alphabet coincides with the permitted prefix character set, so every rendered
// digit is a legal prefix symbol and needs no separate escaping.
type Base36Codec struct {
	base36Alphabet *Alphabet
}

// NewBase36Codec builds a Base36Codec bound to the given Base36 alphabet.
func NewBase36Codec(base36Alphabet *Alphabet) *Base36Codec {
	return &Base36Codec{base36Alphabet: base36Alphabet}
}

// Encode renders the given value as the shortest Base36 string, emitting a single
// zero digit when the value is zero. The modulo loop yields least-significant digits
// first, so the accumulated digits are reversed into big-endian order before return.
func (base36Codec *Base36Codec) Encode(value uint64) string {
	if value == 0 {
		return string(base36Codec.base36Alphabet.CharacterAt(0))
	}

	var radix uint64 = uint64(base36Codec.base36Alphabet.Size())
	var digits []byte = make([]byte, 0, 13)
	var remaining uint64 = value

	for remaining > 0 {
		var remainder uint64 = remaining % radix
		digits = append(digits, base36Codec.base36Alphabet.CharacterAt(int(remainder)))
		remaining = remaining / radix
	}

	var left int = 0
	var right int = len(digits) - 1

	for left < right {
		digits[left], digits[right] = digits[right], digits[left]
		left++
		right--
	}

	return string(digits)
}

// Decode parses a Base36 string back into its unsigned integer value, rejecting an
// empty input, any character outside the alphabet, or a magnitude that overflows a
// 64-bit unsigned integer.
func (base36Codec *Base36Codec) Decode(text string) (uint64, error) {
	if len(text) == 0 {
		return 0, errors.New("base36 codec cannot decode an empty string")
	}

	var radix uint64 = uint64(base36Codec.base36Alphabet.Size())
	var value uint64 = 0

	var position int

	for position = 0; position < len(text); position++ {
		var index int
		var present bool
		index, present = base36Codec.base36Alphabet.IndexOf(text[position])

		if !present {
			return 0, errors.New("base36 codec encountered a character outside the alphabet")
		}

		if value > (math.MaxUint64-uint64(index))/radix {
			return 0, errors.New("base36 codec input overflows a 64-bit unsigned integer")
		}

		value = value*radix + uint64(index)
	}

	return value, nil
}
