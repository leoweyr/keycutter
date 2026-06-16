import { Crc32Calculator } from "./checksum/Crc32Calculator";
import { Alphabet } from "./encoding/Alphabet";
import { Base36Codec } from "./encoding/Base36Codec";
import { Base62Codec } from "./encoding/Base62Codec";
import { SecureGenerator } from "./entropy/SecureGenerator";
import { Clock } from "./timestamp/Clock";
import { SystemClock } from "./timestamp/SystemClock";
import { TokenGenerator } from "./token/TokenGenerator";
import { Token } from "./token/Token";
import { TokenValidator } from "./token/TokenValidator";
import { SecurityContext } from "./SecurityContext";


/**
 * The public entry point of the Tokenforge module. It composes the generation and validation
 * pipelines and exposes them behind a small surface.
 */
export class Forge {
    private readonly _generator: TokenGenerator;
    private readonly _validator: TokenValidator;

    public constructor() {
        const base62Alphabet: Alphabet = new Alphabet(Alphabet.BASE62_CHARACTERS);
        const prefixAlphabet: Alphabet = new Alphabet(Alphabet.BASE36_CHARACTERS);

        const encoder: Base62Codec = new Base62Codec(base62Alphabet);
        const checksumCalculator: Crc32Calculator = new Crc32Calculator(encoder);
        const entropyGenerator: SecureGenerator = new SecureGenerator(base62Alphabet);

        // The prefix alphabet is exactly the 36-symbol Base36 dictionary in ascending order, so it
        // doubles as the codec alphabet that keeps every timestamp digit a legal prefix symbol.
        const base36Codec: Base36Codec = new Base36Codec(prefixAlphabet);
        const clock: Clock = SystemClock.getInstance();

        this._generator = new TokenGenerator(prefixAlphabet, entropyGenerator, checksumCalculator, clock, base36Codec);
        this._validator = new TokenValidator(prefixAlphabet, checksumCalculator, base36Codec);
    }

    /**
     * Runs the full generation pipeline for the given semantic identifiers and returns the rendered
     * token string.
     */
    public generate(systemIdentifier: string, environmentIdentifier: string, domainPurposeIdentifier: string): string {
        const generated: Token = this._generator.generate(systemIdentifier, environmentIdentifier, domainPurposeIdentifier);

        return generated.toString();
    }

    /**
     * Runs the generation pipeline for the given semantic identifiers while appending the current
     * instant as an optional Base36 Unix-seconds timestamp, and returns the rendered token string.
     */
    public generateWithTimestamp(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
    ): string {
        const generated: Token = this._generator.generateWithTimestamp(
            systemIdentifier,
            environmentIdentifier,
            domainPurposeIdentifier,
        );

        return generated.toString();
    }

    /**
     * Runs the full validation pipeline against the raw token and returns the reified security
     * context when the token is structurally and cryptographically sound. The validator detects an
     * embedded timestamp from the prefix component count alone, without any prior knowledge of
     * whether the token carries one.
     */
    public validate(rawToken: string): SecurityContext {
        const validated: Token = this._validator.validate(rawToken);

        return new SecurityContext(
            validated.getSystemIdentifier(),
            validated.getEnvironmentIdentifier(),
            validated.getDomainPurposeIdentifier(),
            validated.getTimestamp(),
        );
    }
}
