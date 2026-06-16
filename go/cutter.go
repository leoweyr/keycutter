package keycutter

import (
	"go.leoweyr.com/keycutter/go/v2/internal/checksum"
	"go.leoweyr.com/keycutter/go/v2/internal/encoding"
	"go.leoweyr.com/keycutter/go/v2/internal/entropy"
	"go.leoweyr.com/keycutter/go/v2/internal/timestamp"
	"go.leoweyr.com/keycutter/go/v2/internal/token"
)

// Cutter is the public entry point of the module. It composes the
// generation and validation pipelines and exposes them behind a small surface.
type Cutter struct {
	generator *token.TokenGenerator
	validator *token.TokenValidator
}

// NewCutter builds a Cutter by wiring the Base62 alphabet, the prefix alphabet, the
// Base62 encoder, the CRC32 checksum calculator, and the secure entropy generator
// into a generation pipeline and a validation pipeline.
func NewCutter() *Cutter {
	var base62Alphabet *encoding.Alphabet = encoding.NewAlphabet(encoding.Base62Characters)
	var prefixAlphabet *encoding.Alphabet = encoding.NewAlphabet(encoding.Base36Characters)

	var encoder *encoding.Base62Codec = encoding.NewBase62Codec(base62Alphabet)
	var checksumCalculator *checksum.Crc32Calculator = checksum.NewCrc32Calculator(encoder)
	var entropyGenerator *entropy.SecureGenerator = entropy.NewSecureGenerator(base62Alphabet)

	// The prefix alphabet is exactly the 36-symbol Base36 dictionary in ascending order,
	// so it doubles as the codec alphabet that keeps every timestamp digit a legal prefix symbol.
	var base36Codec *encoding.Base36Codec = encoding.NewBase36Codec(prefixAlphabet)
	var clock *timestamp.SystemClock = timestamp.SharedSystemClock()

	var generator *token.TokenGenerator = token.NewTokenGenerator(prefixAlphabet, entropyGenerator, checksumCalculator, clock, base36Codec)
	var validator *token.TokenValidator = token.NewTokenValidator(prefixAlphabet, checksumCalculator, base36Codec)

	return &Cutter{
		generator: generator,
		validator: validator,
	}
}

// Generate runs the full generation pipeline for the given semantic identifiers and
// returns the rendered token string.
func (cutter *Cutter) Generate(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (string, error) {
	var generated *token.Token
	var generationError error
	generated, generationError = cutter.generator.Generate(systemIdentifier, environmentIdentifier, domainPurposeIdentifier)

	if generationError != nil {
		return "", generationError
	}

	return generated.String(), nil
}

// GenerateWithTimestamp runs the generation pipeline for the given semantic identifiers
// while appending the current instant as an optional Base36 Unix-seconds timestamp, and
// returns the rendered token string.
func (cutter *Cutter) GenerateWithTimestamp(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (string, error) {
	var generated *token.Token
	var generationError error
	generated, generationError = cutter.generator.GenerateWithTimestamp(systemIdentifier, environmentIdentifier, domainPurposeIdentifier)

	if generationError != nil {
		return "", generationError
	}

	return generated.String(), nil
}

// Validate runs the full validation pipeline against the raw token and returns the
// reified security context when the token is structurally and cryptographically sound.
// The validator detects an embedded timestamp from the prefix component count alone,
// without any prior knowledge of whether the token carries one.
func (cutter *Cutter) Validate(rawToken string) (*SecurityContext, error) {
	var validated *token.Token
	var validationError error
	validated, validationError = cutter.validator.Validate(rawToken)

	if validationError != nil {
		return nil, validationError
	}

	return newSecurityContext(validated.SystemIdentifier(), validated.EnvironmentIdentifier(), validated.DomainPurposeIdentifier(), validated.Timestamp()), nil
}
