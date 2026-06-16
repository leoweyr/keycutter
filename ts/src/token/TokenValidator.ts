import { Calculator } from "../checksum/Calculator";
import { Checksum } from "../checksum/Checksum";
import { Alphabet } from "../encoding/Alphabet";
import { Base36Codec } from "../encoding/Base36Codec";
import { Entropy } from "../entropy/Entropy";
import { ValidationError } from "../errors/ValidationError";
import { Timestamp } from "../timestamp/Timestamp";
import { Token } from "./Token";
import { ReifiedContext } from "./ReifiedContext";


/**
 * Orchestrates the validation pipeline: a structural guard, an asymmetric checksum
 * slice, idempotent integrity verification, and context reification into the token's three prefix
 * identifiers plus an optional timestamp.
 */
export class TokenValidator {
    private readonly _prefixAlphabet: Alphabet;
    private readonly _checksumCalculator: Calculator;
    private readonly _base36Codec: Base36Codec;
    
    public constructor(prefixAlphabet: Alphabet, checksumCalculator: Calculator, base36Codec: Base36Codec) {
        this._prefixAlphabet = prefixAlphabet;
        this._checksumCalculator = checksumCalculator;
        this._base36Codec = base36Codec;
    }

    /**
     * Re-derives the expected checksum from the base string and rejects the token when it diverges
     * from the provided checksum.
     */
    private verifyChecksum(baseString: string, provided: Checksum): void {
        const expected: Checksum = this._checksumCalculator.calculate(baseString);

        if (!expected.equals(provided)) {
            throw new ValidationError("Token checksum does not match its base string.");
        }
    }

    /**
     * Rejects an empty identifier or one bearing a character outside the permitted prefix alphabet.
     */
    private validateComponent(componentName: string, value: string): void {
        if (value.length === 0) {
            throw new ValidationError("Token prefix " + componentName + " component is empty.");
        }

        if (!this._prefixAlphabet.permits(value)) {
            throw new ValidationError(
                "Token prefix " + componentName + " component contains a character outside the permitted set.",
            );
        }
    }

    /**
     * Validates the optional fourth prefix component against the prefix alphabet and decodes its
     * Base36 digits into a Unix-seconds timestamp.
     */
    private reifyTimestamp(component: string): Timestamp {
        this.validateComponent("timestamp", component);

        let seconds: bigint;

        try {
            seconds = this._base36Codec.decode(component);
        } catch {
            throw new ValidationError("Token prefix timestamp component is not a valid Base36 value.");
        }

        return new Timestamp(seconds, component);
    }

    /**
     * Splits the prefix portion into its three semantic identifiers and an optional Base36 timestamp.
     * Because prefix components forbid the separator, the component count alone determines presence:
     * three components carry no timestamp, four carry one as the trailing component, and any other
     * count is rejected.
     */
    private reifyContext(prefixPortion: string): ReifiedContext {
        if (!prefixPortion.endsWith(Token.SEPARATOR)) {
            throw new ValidationError("Token prefix is not terminated by a separator.");
        }

        const core: string = prefixPortion.slice(0, prefixPortion.length - Token.SEPARATOR.length);
        const components: string[] = core.split(Token.SEPARATOR);

        if (components.length !== Token.PREFIX_COMPONENT_COUNT && components.length !== Token.PREFIX_COMPONENT_COUNT + 1) {
            throw new ValidationError("Token prefix does not contain three semantic components with an optional timestamp.");
        }

        this.validateComponent("system", components[0]!);
        this.validateComponent("environment", components[1]!);
        this.validateComponent("domain purpose", components[2]!);

        let reifiedTimestamp: Timestamp | null = null;

        if (components.length === Token.PREFIX_COMPONENT_COUNT + 1) {
            reifiedTimestamp = this.reifyTimestamp(components[Token.PREFIX_COMPONENT_COUNT]!);
        }

        return {
            systemIdentifier: components[0]!,
            environmentIdentifier: components[1]!,
            domainPurposeIdentifier: components[2]!,
            timestamp: reifiedTimestamp,
        };
    }

    /**
     * Enforces the structural guard, performs the asymmetric checksum slice, verifies integrity
     * idempotently, and reifies the semantic context into a token.
     */
    public validate(rawToken: string): Token {
        if (rawToken.length < Token.MINIMUM_LENGTH) {
            throw new ValidationError("Token is shorter than the minimum derived length.");
        }

        const baseString: string = rawToken.slice(0, rawToken.length - Checksum.LENGTH);
        const providedChecksum: Checksum = new Checksum(rawToken.slice(rawToken.length - Checksum.LENGTH));

        this.verifyChecksum(baseString, providedChecksum);

        const prefixPortion: string = baseString.slice(0, baseString.length - Entropy.LENGTH);
        const entropyValue: string = baseString.slice(baseString.length - Entropy.LENGTH);

        const context: ReifiedContext = this.reifyContext(prefixPortion);

        return new Token(
            context.systemIdentifier,
            context.environmentIdentifier,
            context.domainPurposeIdentifier,
            context.timestamp,
            new Entropy(entropyValue),
            providedChecksum,
        );
    }
}
