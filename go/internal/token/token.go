package token

import (
	"go.leoweyr.com/tokenforge/go/internal/checksum"
	"go.leoweyr.com/tokenforge/go/internal/entropy"
	"go.leoweyr.com/tokenforge/go/internal/timestamp"
)

// Separator is the single-byte delimiter placed between prefix semantic components
// and between the prefix and the entropy segment.
const Separator string = "_"

// PrefixComponentCount is the exact number of semantic identifiers required inside
// a well-formed prefix.
const PrefixComponentCount int = 3

// MinimumPrefixLength is the shortest legal prefix portion, reached when each of the
// three semantic components holds a single character (Plus the three delimiters).
const MinimumPrefixLength int = PrefixComponentCount * 2

// MinimumLength is the shortest legal token, derived from the credential topology
// rather than any hard-coded magic number.
const MinimumLength int = MinimumPrefixLength + entropy.Length + checksum.Length

// assembleBaseString concatenates the three prefix identifiers, the optional Base36
// timestamp component, the separators, and the entropy segment into the canonical
// base string used for checksum mapping. A nil timestamp yields the bare three-component
// prefix; a present timestamp appends a fourth component after the domain purpose.
func assembleBaseString(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string, timestamp *timestamp.Timestamp, entropy *entropy.Entropy) string {
	var prefix string = systemIdentifier + Separator + environmentIdentifier + Separator + domainPurposeIdentifier + Separator

	if timestamp != nil {
		prefix = prefix + timestamp.Value() + Separator
	}

	return prefix + entropy.Value()
}

// Token is a fully assembled credential storing the three prefix semantic identifiers,
// an optional Base36 creation timestamp, a high-intensity entropy segment, and a tail
// checksum.
type Token struct {
	systemIdentifier        string
	environmentIdentifier   string
	domainPurposeIdentifier string
	timestamp               *timestamp.Timestamp
	entropy                 *entropy.Entropy
	checksum                *checksum.Checksum
}

// NewToken composes a Token from its three prefix identifiers, an optional timestamp,
// entropy, and checksum. A nil timestamp marks a token without the optional component.
func NewToken(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string, timestamp *timestamp.Timestamp, entropy *entropy.Entropy, checksum *checksum.Checksum) *Token {
	return &Token{
		systemIdentifier:        systemIdentifier,
		environmentIdentifier:   environmentIdentifier,
		domainPurposeIdentifier: domainPurposeIdentifier,
		timestamp:               timestamp,
		entropy:                 entropy,
		checksum:                checksum,
	}
}

// SystemIdentifier returns the system identifier.
func (token *Token) SystemIdentifier() string {
	return token.systemIdentifier
}

// EnvironmentIdentifier returns the environment identifier.
func (token *Token) EnvironmentIdentifier() string {
	return token.environmentIdentifier
}

// DomainPurposeIdentifier returns the domain purpose identifier.
func (token *Token) DomainPurposeIdentifier() string {
	return token.domainPurposeIdentifier
}

// Timestamp returns the optional creation timestamp, or nil when the token carries none.
func (token *Token) Timestamp() *timestamp.Timestamp {
	return token.timestamp
}

// Entropy returns the high-intensity entropy segment.
func (token *Token) Entropy() *entropy.Entropy {
	return token.entropy
}

// Checksum returns the tail checksum segment.
func (token *Token) Checksum() *checksum.Checksum {
	return token.checksum
}

// BaseString returns the checksum-free base string (prefix + entropy, joined by separators).
func (token *Token) BaseString() string {
	return assembleBaseString(token.systemIdentifier, token.environmentIdentifier, token.domainPurposeIdentifier, token.timestamp, token.entropy)
}

// String renders the complete token.
func (token *Token) String() string {
	return token.BaseString() + token.checksum.Value()
}
