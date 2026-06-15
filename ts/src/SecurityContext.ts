import { Timestamp } from "./timestamp/Timestamp";
import { Token } from "./token/Token";


/**
 * SecurityContext exposes the plaintext semantic identifiers reified from a validated token — the
 * system, the environment, and the domain purpose — together with the optional creation timestamp
 * when the token carried one.
 */
export class SecurityContext {
    private readonly _systemIdentifier: string;
    private readonly _environmentIdentifier: string;
    private readonly _domainPurposeIdentifier: string;
    private readonly _timestamp: Timestamp | null;

    public constructor(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
        timestamp: Timestamp | null,
    ) {
        this._systemIdentifier = systemIdentifier;
        this._environmentIdentifier = environmentIdentifier;
        this._domainPurposeIdentifier = domainPurposeIdentifier;
        this._timestamp = timestamp;
    }

    /** @returns The system identifier. */
    public getSystemIdentifier(): string {
        return this._systemIdentifier;
    }

    /** @returns The environment identifier. */
    public getEnvironmentIdentifier(): string {
        return this._environmentIdentifier;
    }

    /** @returns The domain purpose identifier. */
    public getDomainPurposeIdentifier(): string {
        return this._domainPurposeIdentifier;
    }

    /**
     * @returns The embedded creation instant as a Unix epoch second count, or null when the token
     * carried no timestamp. The null result disambiguates a genuine epoch-zero instant (0n) from an
     * absent timestamp.
     */
    public createdAtUnixSeconds(): bigint | null {
        if (this._timestamp === null) {
            return null;
        }

        return this._timestamp.getSeconds();
    }

    /**
     * Renders the semantic identifiers joined by underscores, appending the Base36 timestamp
     * component when the token carried one.
     */
    public toString(): string {
        const components: string[] = [this._systemIdentifier, this._environmentIdentifier, this._domainPurposeIdentifier];

        if (this._timestamp !== null) {
            components.push(this._timestamp.getValue());
        }

        return components.join(Token.SEPARATOR);
    }
}
