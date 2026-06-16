package keycutter

import (
	"strings"

	"go.leoweyr.com/keycutter/go/v2/internal/timestamp"
)

// SecurityContext exposes the plaintext semantic identifiers reified from a validated
// token — the system, the environment, and the domain purpose — together with the
// optional creation timestamp when the token carried one.
type SecurityContext struct {
	systemIdentifier        string
	environmentIdentifier   string
	domainPurposeIdentifier string
	timestamp               *timestamp.Timestamp
}

// newSecurityContext assembles a SecurityContext from its three semantic identifiers
// and an optional creation timestamp.
func newSecurityContext(systemIdentifier string, environmentIdentifier string, domainIdentifier string, reifiedTimestamp *timestamp.Timestamp) *SecurityContext {
	return &SecurityContext{
		systemIdentifier:        systemIdentifier,
		environmentIdentifier:   environmentIdentifier,
		domainPurposeIdentifier: domainIdentifier,
		timestamp:               reifiedTimestamp,
	}
}

// SystemIdentifier returns the system identifier.
func (securityContext *SecurityContext) SystemIdentifier() string {
	return securityContext.systemIdentifier
}

// EnvironmentIdentifier returns the environment identifier.
func (securityContext *SecurityContext) EnvironmentIdentifier() string {
	return securityContext.environmentIdentifier
}

// DomainPurposeIdentifier returns the domain purpose identifier.
func (securityContext *SecurityContext) DomainPurposeIdentifier() string {
	return securityContext.domainPurposeIdentifier
}

// CreatedAtUnixSeconds returns the embedded creation instant as a Unix epoch second
// count and reports whether the token carried a timestamp. The boolean disambiguates a
// genuine epoch-zero instant from an absent timestamp.
func (securityContext *SecurityContext) CreatedAtUnixSeconds() (uint64, bool) {
	if securityContext.timestamp == nil {
		return 0, false
	}

	return securityContext.timestamp.Seconds(), true
}

// String renders the semantic identifiers joined by underscores, appending the Base36
// timestamp component when the token carried one.
func (securityContext *SecurityContext) String() string {
	var components []string = []string{securityContext.systemIdentifier, securityContext.environmentIdentifier, securityContext.domainPurposeIdentifier}

	if securityContext.timestamp != nil {
		components = append(components, securityContext.timestamp.Value())
	}

	return strings.Join(components, "_")
}
