package encoding

// Base62Encoder converts unsigned integers into fixed-width Base62 strings using a
// right-to-left modulo loop that yields natural front zero-padding.
type Base62Encoder struct {
	base62Alphabet *Alphabet
}

// NewBase62Encoder builds a Base62Encoder bound to the given Base62 alphabet.
func NewBase62Encoder(base62Alphabet *Alphabet) *Base62Encoder {
	return &Base62Encoder{base62Alphabet: base62Alphabet}
}

// Encode renders the given value as a Base62 string of exactly the requested width,
// filling positions from the least significant digit backward.
func (base62Encoder *Base62Encoder) Encode(value uint32, width int) string {
	var radix uint32 = uint32(base62Encoder.base62Alphabet.Size())
	var characters []byte = make([]byte, width)
	var remaining uint32 = value

	var position int

	for position = width - 1; position >= 0; position-- {
		var remainder uint32 = remaining % radix
		characters[position] = base62Encoder.base62Alphabet.CharacterAt(int(remainder))
		remaining = remaining / radix
	}

	return string(characters)
}
