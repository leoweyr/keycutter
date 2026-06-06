package token

import (
	"go.leoweyr.com/tokenforge/go/internal/checksum"
	"go.leoweyr.com/tokenforge/go/internal/encoding"
	"go.leoweyr.com/tokenforge/go/internal/entropy"
	"go.leoweyr.com/tokenforge/go/internal/fault"
)

// TokenGenerator orchestrates the full generation pipeline, turning a set of
// semantic identifiers into a verified, fully assembled token.
type TokenGenerator struct {
	prefixAlphabet     *encoding.Alphabet
	entropyGenerator   entropy.Generator
	checksumCalculator checksum.Calculator
}

// NewTokenGenerator wires a TokenGenerator to its prefix alphabet, entropy source,
// and checksum calculator.
func NewTokenGenerator(prefixAlphabet *encoding.Alphabet, entropyGenerator entropy.Generator, checksumCalculator checksum.Calculator) *TokenGenerator {
	return &TokenGenerator{
		prefixAlphabet:     prefixAlphabet,
		entropyGenerator:   entropyGenerator,
		checksumCalculator: checksumCalculator,
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

// Generate validates the supplied identifiers, draws unbiased entropy, fuses the
// base string, maps the checksum, and returns the assembled token.
func (tokenGenerator *TokenGenerator) Generate(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (*Token, error) {
	var systemError error = tokenGenerator.validateComponent("System", systemIdentifier)

	if systemError != nil {
		return nil, systemError
	}

	var environmentError error = tokenGenerator.validateComponent("Environment", environmentIdentifier)

	if environmentError != nil {
		return nil, environmentError
	}

	var domainError error = tokenGenerator.validateComponent("Domain purpose", domainPurposeIdentifier)

	if domainError != nil {
		return nil, domainError
	}

	var entropySegment *entropy.Entropy
	var generationError error
	entropySegment, generationError = tokenGenerator.entropyGenerator.Generate()

	if generationError != nil {
		return nil, generationError
	}

	var baseString string = assembleBaseString(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, entropySegment)
	var checksumSegment *checksum.Checksum = tokenGenerator.checksumCalculator.Calculate(baseString)

	return NewToken(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, entropySegment, checksumSegment), nil
}
