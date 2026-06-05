package token

import (
	"strings"

	"go.leoweyr.com/tokenforge/go/internal/checksum"
	"go.leoweyr.com/tokenforge/go/internal/encoding"
	"go.leoweyr.com/tokenforge/go/internal/entropy"
	"go.leoweyr.com/tokenforge/go/internal/fault"
)

// TokenValidator orchestrates the validation pipeline: a structural guard, an
// asymmetric checksum slice, idempotent integrity verification, and context
// reification into the token's three prefix identifiers.
type TokenValidator struct {
	prefixAlphabet     *encoding.Alphabet
	checksumCalculator checksum.Calculator
}

// NewTokenValidator wires a TokenValidator to its prefix alphabet and checksum calculator.
func NewTokenValidator(prefixAlphabet *encoding.Alphabet, checksumCalculator checksum.Calculator) *TokenValidator {
	return &TokenValidator{
		prefixAlphabet:     prefixAlphabet,
		checksumCalculator: checksumCalculator,
	}
}

// verifyChecksum re-derives the expected checksum from the base string and rejects
// the token when it diverges from the provided checksum.
func (tokenValidator *TokenValidator) verifyChecksum(baseString string, provided *checksum.Checksum) error {
	var expected *checksum.Checksum = tokenValidator.checksumCalculator.Calculate(baseString)

	if !expected.Equals(provided) {
		return fault.NewValidationError("Token checksum does not match its base string")
	}

	return nil
}

// validateComponent rejects an empty identifier or one bearing a character outside
// the permitted prefix alphabet.
func (tokenValidator *TokenValidator) validateComponent(componentName string, value string) error {
	if len(value) == 0 {
		return fault.NewValidationError("Token prefix " + componentName + " component is empty")
	}

	if !tokenValidator.prefixAlphabet.Permits(value) {
		return fault.NewValidationError("Token prefix " + componentName + " component contains a character outside the permitted set")
	}

	return nil
}

// reifyContext splits the prefix portion into its three semantic identifiers,
// enforcing the exact component count and the permitted prefix alphabet.
func (tokenValidator *TokenValidator) reifyContext(prefixPortion string) (string, string, string, error) {
	if !strings.HasSuffix(prefixPortion, Separator) {
		return "", "", "", fault.NewValidationError("Token prefix is not terminated by a separator")
	}

	var core string = prefixPortion[:len(prefixPortion)-len(Separator)]
	var components []string = strings.Split(core, Separator)

	if len(components) != PrefixComponentCount {
		return "", "", "", fault.NewValidationError("Token prefix does not contain exactly three semantic components")
	}

	var systemError error = tokenValidator.validateComponent("system", components[0])

	if systemError != nil {
		return "", "", "", systemError
	}

	var environmentError error = tokenValidator.validateComponent("environment", components[1])

	if environmentError != nil {
		return "", "", "", environmentError
	}

	var domainPurposeError error = tokenValidator.validateComponent("domain purpose", components[2])

	if domainPurposeError != nil {
		return "", "", "", domainPurposeError
	}

	return components[0], components[1], components[2], nil
}

// Validate enforces the structural guard, performs the asymmetric checksum slice,
// verifies integrity idempotently, and reifies the semantic context into a token.
func (tokenValidator *TokenValidator) Validate(rawToken string) (*Token, error) {
	if len(rawToken) < MinimumLength {
		return nil, fault.NewValidationError("Token is shorter than the minimum derived length")
	}

	var baseString string = rawToken[:len(rawToken)-checksum.Length]
	var providedChecksum *checksum.Checksum = checksum.NewChecksum(rawToken[len(rawToken)-checksum.Length:])

	var verificationError error = tokenValidator.verifyChecksum(baseString, providedChecksum)

	if verificationError != nil {
		return nil, verificationError
	}

	var prefixPortion string = baseString[:len(baseString)-entropy.Length]
	var entropyValue string = baseString[len(baseString)-entropy.Length:]

	var system string
	var environment string
	var domain string
	var reificationError error
	system, environment, domain, reificationError = tokenValidator.reifyContext(prefixPortion)

	if reificationError != nil {
		return nil, reificationError
	}

	return NewToken(system, environment, domain, entropy.NewEntropy(entropyValue), providedChecksum), nil
}
