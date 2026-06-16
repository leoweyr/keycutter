package token

import (
	"go.leoweyr.com/tokenforge/go/internal/checksum"
	"go.leoweyr.com/tokenforge/go/internal/encoding"
	"go.leoweyr.com/tokenforge/go/internal/entropy"
	"go.leoweyr.com/tokenforge/go/internal/fault"
	"go.leoweyr.com/tokenforge/go/internal/timestamp"
)

// TokenGenerator orchestrates the full generation pipeline, turning a set of
// semantic identifiers into a verified, fully assembled token.
type TokenGenerator struct {
	prefixAlphabet     *encoding.Alphabet
	entropyGenerator   entropy.Generator
	checksumCalculator checksum.Calculator
	clock              timestamp.Clock
	base36Codec        *encoding.Base36Codec
}

// NewTokenGenerator wires a TokenGenerator to its prefix alphabet, entropy source,
// checksum calculator, wall clock, and Base36 codec.
func NewTokenGenerator(prefixAlphabet *encoding.Alphabet, entropyGenerator entropy.Generator, checksumCalculator checksum.Calculator, clock timestamp.Clock, base36Codec *encoding.Base36Codec) *TokenGenerator {
	return &TokenGenerator{
		prefixAlphabet:     prefixAlphabet,
		entropyGenerator:   entropyGenerator,
		checksumCalculator: checksumCalculator,
		clock:              clock,
		base36Codec:        base36Codec,
	}
}

// validateComponent rejects an empty semantic identifier or one carrying any
// character outside the permitted prefix alphabet.
func (tokenGenerator *TokenGenerator) validateComponent(label string, value string) error {
	if len(value) == 0 {
		return fault.NewValidationError(label + " identifier must not be empty")
	}

	if !tokenGenerator.prefixAlphabet.Permits(value) {
		return fault.NewValidationError(label + " identifier contains a character outside the permitted set")
	}

	return nil
}

// Validates the supplied identifiers, draws unbiased entropy, fuses the
// base string around the optional timestamp, maps the checksum, and returns the
// assembled token.
func (tokenGenerator *TokenGenerator) build(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string, tokenTimestamp *timestamp.Timestamp) (*Token, error) {
	var systemError error = tokenGenerator.validateComponent("system", systemIdentifier)

	if systemError != nil {
		return nil, systemError
	}

	var environmentError error = tokenGenerator.validateComponent("environment", environmentIdentifier)

	if environmentError != nil {
		return nil, environmentError
	}

	var domainError error = tokenGenerator.validateComponent("domain purpose", domainPurposeIdentifier)

	if domainError != nil {
		return nil, domainError
	}

	var entropySegment *entropy.Entropy
	var generationError error
	entropySegment, generationError = tokenGenerator.entropyGenerator.Generate()

	if generationError != nil {
		return nil, generationError
	}

	var baseString string = assembleBaseString(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, tokenTimestamp, entropySegment)
	var checksumSegment *checksum.Checksum = tokenGenerator.checksumCalculator.Calculate(baseString)

	return NewToken(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, tokenTimestamp, entropySegment, checksumSegment), nil
}

// Generate produces a token whose prefix carries only the three semantic identifiers.
func (tokenGenerator *TokenGenerator) Generate(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (*Token, error) {
	return tokenGenerator.build(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, nil)
}

// GenerateWithTimestamp produces a token whose prefix appends a fourth component: the
// current wall-clock instant as Unix seconds rendered in Base36, marking the credential
// version in a self-describing way for distributed systems.
func (tokenGenerator *TokenGenerator) GenerateWithTimestamp(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (*Token, error) {
	var seconds uint64 = tokenGenerator.clock.NowUnixSeconds()
	var encoded string = tokenGenerator.base36Codec.Encode(seconds)
	var tokenCreated *timestamp.Timestamp = timestamp.NewTimestamp(seconds, encoded)

	return tokenGenerator.build(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, tokenCreated)
}
