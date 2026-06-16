/**
 * ValidationError represents a structural or integrity failure detected while a token is checked
 * against the Tokenforge specification.
 */
export class ValidationError extends Error {
    /**
     * Constructs a ValidationError carrying the given human-readable reason.
     */
    public constructor(reason: string) {
        super(reason);
        this.name = "ValidationError";

        // Restore the prototype chain so instanceof works after transpilation to ES5 targets.
        Object.setPrototypeOf(this, ValidationError.prototype);
    }
}
