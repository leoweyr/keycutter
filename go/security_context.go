package tokenforge

import "strings"

// SecurityContext exposes the three plaintext semantic identifiers reified from a
// validated token: the system, the environment, and the domain purpose.
type SecurityContext struct {
	systemIdentifier        string
	environmentIdentifier   string
	domainPurposeIdentifier string
}

// newSecurityContext assembles a SecurityContext from its three semantic identifiers.
func newSecurityContext(systemIdentifier string, environmentIdentifier string, domainIdentifier string) *SecurityContext {
	return &SecurityContext{
		systemIdentifier:        systemIdentifier,
		environmentIdentifier:   environmentIdentifier,
		domainPurposeIdentifier: domainIdentifier,
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

// String renders the three identifiers joined by underscores.
func (securityContext *SecurityContext) String() string {
	return strings.Join([]string{securityContext.systemIdentifier, securityContext.environmentIdentifier, securityContext.domainPurposeIdentifier}, "_")
}
