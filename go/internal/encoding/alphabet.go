package encoding

// Base62Characters is the canonical 62-symbol GMP dictionary ordered by ascending
// ASCII code, used for high-intensity entropy and tail checksum encoding.
const Base62Characters string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base36Characters is the canonical 36-symbol dictionary ordered by ascending ASCII
// code, spanning the decimal digits and lowercase ASCII letters.
const Base36Characters string = "0123456789abcdefghijklmnopqrstuvwxyz"

// Alphabet is an immutable ordered symbol set that maps between character values
// and their positional indices.
type Alphabet struct {
	characters       string
	indexByCharacter map[byte]int
}

// NewAlphabet builds an Alphabet from the given ordered character set and indexes
// every symbol for constant-time membership and lookup queries.
func NewAlphabet(characters string) *Alphabet {
	var indexByCharacter map[byte]int = make(map[byte]int, len(characters))

	var position int

	for position = 0; position < len(characters); position++ {
		indexByCharacter[characters[position]] = position
	}

	return &Alphabet{
		characters:       characters,
		indexByCharacter: indexByCharacter,
	}
}

// Size returns the number of symbols contained in the alphabet.
func (alphabet *Alphabet) Size() int {
	return len(alphabet.characters)
}

// CharacterAt returns the symbol located at the given positional index.
func (alphabet *Alphabet) CharacterAt(index int) byte {
	return alphabet.characters[index]
}

// Contains reports whether the given character belongs to the alphabet.
func (alphabet *Alphabet) Contains(character byte) bool {
	var present bool
	_, present = alphabet.indexByCharacter[character]

	return present
}

// IndexOf returns the positional index of the given character and whether that
// character belongs to the alphabet.
func (alphabet *Alphabet) IndexOf(character byte) (int, bool) {
	var index int
	var present bool
	index, present = alphabet.indexByCharacter[character]

	return index, present
}

// Permits reports whether every character in the text belongs to the alphabet.
func (alphabet *Alphabet) Permits(text string) bool {
	var position int

	for position = 0; position < len(text); position++ {
		if !alphabet.Contains(text[position]) {
			return false
		}
	}

	return true
}
