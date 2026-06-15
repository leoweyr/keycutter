import { Alphabet } from "./Alphabet";


/**
 * Base36Codec converts unsigned integers to and from variable-width Base36 strings. Its alphabet
 * coincides with the permitted prefix character set, so every rendered digit is a legal prefix
 * symbol and needs no separate escaping. Values are carried as bigint so the full unsigned 64-bit
 * domain survives without precision loss.
 */
export class Base36Codec {
    /**
     * The largest value representable by a 64-bit unsigned integer, mirroring Go's math.MaxUint64
     * so decode rejects the exact same overflow set.
     */
    private static readonly _MAXIMUM_UNSIGNED_SIXTY_FOUR_BIT: bigint = (1n << 64n) - 1n;

    private readonly _base36Alphabet: Alphabet;

    /**
     * Constructs a Base36Codec bound to the given Base36 alphabet.
     */
    public constructor(base36Alphabet: Alphabet) {
        this._base36Alphabet = base36Alphabet;
    }

    /**
     * Renders the given value as the shortest Base36 string, emitting a single zero digit when the
     * value is zero. The modulo loop yields least-significant digits first, so the accumulated
     * digits are reversed into big-endian order before return.
     */
    public encode(value: bigint): string {
        if (value === 0n) {
            return this._base36Alphabet.characterAt(0);
        }

        const radix: bigint = BigInt(this._base36Alphabet.size());
        const digits: string[] = [];
        let remaining: bigint = value;

        while (remaining > 0n) {
            const remainder: bigint = remaining % radix;
            digits.push(this._base36Alphabet.characterAt(Number(remainder)));
            remaining = remaining / radix;
        }

        let left: number = 0;
        let right: number = digits.length - 1;

        while (left < right) {
            const swap: string = digits[left]!;
            digits[left] = digits[right]!;
            digits[right] = swap;
            left++;
            right--;
        }

        return digits.join("");
    }

    /**
     * Parses a Base36 string back into its unsigned integer value, rejecting an empty input, any
     * character outside the alphabet, or a magnitude that overflows a 64-bit unsigned integer.
     */
    public decode(text: string): bigint {
        if (text.length === 0) {
            throw new Error("Base36 codec cannot decode an empty string.");
        }

        const radix: bigint = BigInt(this._base36Alphabet.size());
        let value: bigint = 0n;

        for (let position: number = 0; position < text.length; position++) {
            const index: number | undefined = this._base36Alphabet.indexOf(text.charCodeAt(position));

            if (index === undefined) {
                throw new Error("Base36 codec encountered a character outside the alphabet.");
            }

            if (value > (Base36Codec._MAXIMUM_UNSIGNED_SIXTY_FOUR_BIT - BigInt(index)) / radix) {
                throw new Error("Base36 codec input overflows a 64-bit unsigned integer.");
            }

            value = value * radix + BigInt(index);
        }

        return value;
    }
}
