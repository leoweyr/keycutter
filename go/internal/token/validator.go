package token

import (
	"strings"

	"go.leoweyr.com/tokenforge/go/internal/checksum"
	"go.leoweyr.com/tokenforge/go/internal/encoding"
	"go.leoweyr.com/tokenforge/go/internal/entropy"
	"go.leoweyr.com/tokenforge/go/internal/fault"
	"go.leoweyr.com/tokenforge/go/internal/timestamp"
)

// TokenValidator orchestrates the validation pipeline: a structural guard, an
// asymmetric checksum slice, idempotent integrity verification, and context
// reification into the token's three prefix identifiers plus an optional timestamp.
type TokenValidator struct {
	prefixAlphabet     *encoding.Alphabet
	checksumCalculator checksum.Calculator
	base36Codec        *encoding.Base36Codec
}

// NewTokenValidator wires a TokenValidator to its prefix alphabet, checksum calculator,
// and Base36 codec.
func NewTokenValidator(prefixAlphabet *encoding.Alphabet, checksumCalculator checksum.Calculator, base36Codec *encoding.Base36Codec) *TokenValidator {
	return &TokenValidator{
		prefixAlphabet:     prefixAlphabet,
		checksumCalculator: checksumCalculator,
		base36Codec:        base36Codec,
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

// reifyTimestamp validates the optional fourth prefix component against the prefix
// alphabet and decodes its Base36 digits into a Unix-seconds timestamp.
func (tokenValidator *TokenValidator) reifyTimestamp(component string) (*timestamp.Timestamp, error) {
	var validationError error = tokenValidator.validateComponent("timestamp", component)

	if validationError != nil {
		return nil, validationError
	}

	var seconds uint64
	var decodeError error
	seconds, decodeError = tokenValidator.base36Codec.Decode(component)

	if decodeError != nil {
		return nil, fault.NewValidationError("Token prefix timestamp component is not a valid Base36 value")
	}

	return timestamp.NewTimestamp(seconds, component), nil
}

// reifyContext splits the prefix portion into its three semantic identifiers and an
// optional Base36 timestamp. Because prefix components forbid the separator, the
// component count alone determines presence: three components carry no timestamp,
// four carry one as the trailing component, and any other count is rejected.
func (tokenValidator *TokenValidator) reifyContext(prefixPortion string) (string, string, string, *timestamp.Timestamp, error) {
	if !strings.HasSuffix(prefixPortion, Separator) {
		return "", "", "", nil, fault.NewValidationError("Token prefix is not terminated by a separator")
	}

	var core string = prefixPortion[:len(prefixPortion)-len(Separator)]
	var components []string = strings.Split(core, Separator)

	if len(components) != PrefixComponentCount && len(components) != PrefixComponentCount+1 {
		return "", "", "", nil, fault.NewValidationError("Token prefix does not contain three semantic components with an optional timestamp")
	}

	var systemError error = tokenValidator.validateComponent("system", components[0])

	if systemError != nil {
		return "", "", "", nil, systemError
	}

	var environmentError error = tokenValidator.validateComponent("environment", components[1])

	if environmentError != nil {
		return "", "", "", nil, environmentError
	}

	var domainPurposeError error = tokenValidator.validateComponent("domain purpose", components[2])

	if domainPurposeError != nil {
		return "", "", "", nil, domainPurposeError
	}

	var reifiedTimestamp *timestamp.Timestamp = nil

	if len(components) == PrefixComponentCount+1 {
		var timestampError error
		reifiedTimestamp, timestampError = tokenValidator.reifyTimestamp(components[PrefixComponentCount])

		if timestampError != nil {
			return "", "", "", nil, timestampError
		}
	}

	return components[0], components[1], components[2], reifiedTimestamp, nil
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
	var reifiedTimestamp *timestamp.Timestamp
	var reificationError error
	system, environment, domain, reifiedTimestamp, reificationError = tokenValidator.reifyContext(prefixPortion)

	if reificationError != nil {
		return nil, reificationError
	}

	return NewToken(system, environment, domain, reifiedTimestamp, entropy.NewEntropy(entropyValue), providedChecksum), nil
}
