package tokenforge

import (
	"go.leoweyr.com/tokenforge/internal/checksum"
	"go.leoweyr.com/tokenforge/internal/encoding"
	"go.leoweyr.com/tokenforge/internal/entropy"
	"go.leoweyr.com/tokenforge/internal/token"
)

// Forge is the public entry point of the Tokenforge module. It composes the
// generation and validation pipelines and exposes them behind a small surface.
type Forge struct {
	generator *token.TokenGenerator
	validator *token.TokenValidator
}

// NewForge builds a Forge by wiring the Base62 alphabet, the prefix alphabet, the
// Base62 encoder, the CRC32 checksum calculator, and the secure entropy generator
// into a generation pipeline and a validation pipeline.
func NewForge() *Forge {
	var base62Alphabet *encoding.Alphabet = encoding.NewAlphabet(encoding.Base62Characters)
	var prefixAlphabet *encoding.Alphabet = encoding.NewAlphabet(encoding.PrefixCharacters)

	var encoder *encoding.Base62Encoder = encoding.NewBase62Encoder(base62Alphabet)
	var checksumCalculator *checksum.Crc32Calculator = checksum.NewCrc32Calculator(encoder)
	var entropyGenerator *entropy.SecureGenerator = entropy.NewSecureGenerator(base62Alphabet)

	var generator *token.TokenGenerator = token.NewTokenGenerator(prefixAlphabet, entropyGenerator, checksumCalculator)
	var validator *token.TokenValidator = token.NewTokenValidator(prefixAlphabet, checksumCalculator)

	return &Forge{
		generator: generator,
		validator: validator,
	}
}

// Generate runs the full generation pipeline for the given semantic identifiers and
// returns the rendered token string.
func (forge *Forge) Generate(systemIdentifier string, environmentIdentifier string, domainPurposeIdentifier string) (string, error) {
	var generated *token.Token
	var generationError error
	generated, generationError = forge.generator.Generate(systemIdentifier, environmentIdentifier, domainPurposeIdentifier)

	if generationError != nil {
		return "", generationError
	}

	return generated.String(), nil
}

// Validate runs the full validation pipeline against the raw token and returns the
// reified security context when the token is structurally and cryptographically sound.
func (forge *Forge) Validate(rawToken string) (*SecurityContext, error) {
	var validated *token.Token
	var validationError error
	validated, validationError = forge.validator.Validate(rawToken)

	if validationError != nil {
		return nil, validationError
	}

	return newSecurityContext(validated.SystemIdentifier(), validated.EnvironmentIdentifier(), validated.DomainPurposeIdentifier()), nil
}
