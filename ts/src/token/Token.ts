import { Checksum } from "../checksum/Checksum";
import { Entropy } from "../entropy/Entropy";
import { Timestamp } from "../timestamp/Timestamp";


/**
 * A fully assembled credential storing the three prefix semantic identifiers, an optional
 * Base36 creation timestamp, a high-intensity entropy segment, and a tail checksum.
 */
export class Token {
    /**
     * The single-byte delimiter placed between prefix semantic components and between
     * the prefix and the entropy segment.
     */
    public static readonly SEPARATOR: string = "_";

    /**
     * The exact number of semantic identifiers required inside a well-formed prefix.
     */
    public static readonly PREFIX_COMPONENT_COUNT: number = 3;

    /**
     * The shortest legal prefix portion, reached when each of the three semantic components holds a
     * single character (Plus the three delimiters).
     */
    public static readonly MINIMUM_PREFIX_LENGTH: number = Token.PREFIX_COMPONENT_COUNT * 2;

    /**
     * The shortest legal token, derived from the credential topology rather than any hard-coded
     * magic number.
     */
    public static readonly MINIMUM_LENGTH: number = Token.MINIMUM_PREFIX_LENGTH + Entropy.LENGTH + Checksum.LENGTH;

    /**
     * Concatenates the three prefix identifiers, the optional Base36 timestamp component, the
     * separators, and the entropy segment into the canonical base string used for checksum mapping.
     * A null timestamp yields the bare three-component prefix. A present timestamp appends a fourth
     * component after the domain purpose.
     */
    public static assembleBaseString(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
        timestamp: Timestamp | null,
        entropy: Entropy,
    ): string {
        let prefix: string =
            systemIdentifier + Token.SEPARATOR + environmentIdentifier + Token.SEPARATOR + domainPurposeIdentifier + Token.SEPARATOR;

        if (timestamp !== null) {
            prefix = prefix + timestamp.getValue() + Token.SEPARATOR;
        }

        return prefix + entropy.getValue();
    }

    private readonly _systemIdentifier: string;
    private readonly _environmentIdentifier: string;
    private readonly _domainPurposeIdentifier: string;
    private readonly _timestamp: Timestamp | null;
    private readonly _entropy: Entropy;
    private readonly _checksum: Checksum;

    public constructor(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
        timestamp: Timestamp | null,
        entropy: Entropy,
        checksum: Checksum,
    ) {
        this._systemIdentifier = systemIdentifier;
        this._environmentIdentifier = environmentIdentifier;
        this._domainPurposeIdentifier = domainPurposeIdentifier;
        this._timestamp = timestamp;
        this._entropy = entropy;
        this._checksum = checksum;
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

    /** @returns The optional creation timestamp, or null when the token carries none. */
    public getTimestamp(): Timestamp | null {
        return this._timestamp;
    }

    /** @returns The high-intensity entropy segment. */
    public getEntropy(): Entropy {
        return this._entropy;
    }

    /** @returns The tail checksum segment. */
    public getChecksum(): Checksum {
        return this._checksum;
    }

    /** @returns The checksum-free base string (prefix + entropy, joined by separators). */
    public getBaseString(): string {
        return Token.assembleBaseString(
            this._systemIdentifier,
            this._environmentIdentifier,
            this._domainPurposeIdentifier,
            this._timestamp,
            this._entropy,
        );
    }

    /** Renders the complete token. */
    public toString(): string {
        return this.getBaseString() + this._checksum.getValue();
    }
}
