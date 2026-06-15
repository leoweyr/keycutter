import { Calculator } from "../checksum/Calculator";
import { Checksum } from "../checksum/Checksum";
import { Alphabet } from "../encoding/Alphabet";
import { Base36Codec } from "../encoding/Base36Codec";
import { Entropy } from "../entropy/Entropy";
import { Generator } from "../entropy/Generator";
import { ValidationError } from "../errors/ValidationError";
import { Clock } from "../timestamp/Clock";
import { Timestamp } from "../timestamp/Timestamp";
import { Token } from "./Token";


/**
 * Orchestrates the full generation pipeline, turning a set of semantic identifiers
 * into a verified, fully assembled token.
 */
export class TokenGenerator {
    private readonly _prefixAlphabet: Alphabet;
    private readonly _entropyGenerator: Generator;
    private readonly _checksumCalculator: Calculator;
    private readonly _clock: Clock;
    private readonly _base36Codec: Base36Codec;

    public constructor(
        prefixAlphabet: Alphabet,
        entropyGenerator: Generator,
        checksumCalculator: Calculator,
        clock: Clock,
        base36Codec: Base36Codec,
    ) {
        this._prefixAlphabet = prefixAlphabet;
        this._entropyGenerator = entropyGenerator;
        this._checksumCalculator = checksumCalculator;
        this._clock = clock;
        this._base36Codec = base36Codec;
    }

    /**
     * Rejects an empty semantic identifier or one carrying any character outside the permitted
     * prefix alphabet.
     */
    private validateComponent(label: string, value: string): void {
        if (value.length === 0) {
            throw new ValidationError(label + " identifier must not be empty.");
        }

        if (!this._prefixAlphabet.permits(value)) {
            throw new ValidationError(label + " identifier contains a character outside the permitted set.");
        }
    }

    /**
     * Validates the supplied identifiers, draws unbiased entropy, fuses the base string around the
     * optional timestamp, maps the checksum, and returns the assembled token.
     */
    private build(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
        tokenTimestamp: Timestamp | null,
    ): Token {
        this.validateComponent("System", systemIdentifier);
        this.validateComponent("Environment", environmentIdentifier);
        this.validateComponent("Domain purpose", domainPurposeIdentifier);

        const entropySegment: Entropy = this._entropyGenerator.generate();

        const baseString: string = Token.assembleBaseString(
            systemIdentifier,
            environmentIdentifier,
            domainPurposeIdentifier,
            tokenTimestamp,
            entropySegment,
        );

        const checksumSegment: Checksum = this._checksumCalculator.calculate(baseString);

        return new Token(
            systemIdentifier,
            environmentIdentifier,
            domainPurposeIdentifier,
            tokenTimestamp,
            entropySegment,
            checksumSegment,
        );
    }

    /**
     * Produces a token whose prefix carries only the three semantic identifiers.
     */
    public generate(systemIdentifier: string, environmentIdentifier: string, domainPurposeIdentifier: string): Token {
        return this.build(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, null);
    }

    /**
     * Produces a token whose prefix appends a fourth component: the current wall-clock instant as
     * Unix seconds rendered in Base36, marking the credential version in a self-describing way for
     * distributed systems.
     */
    public generateWithTimestamp(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
    ): Token {
        const seconds: bigint = this._clock.nowUnixSeconds();
        const encoded: string = this._base36Codec.encode(seconds);
        const tokenCreated: Timestamp = new Timestamp(seconds, encoded);

        return this.build(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, tokenCreated);
    }
}
