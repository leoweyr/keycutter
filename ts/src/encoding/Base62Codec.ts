import { Alphabet } from "./Alphabet";


/**
 * Base62Codec converts unsigned integers into fixed-width Base62 strings using a right-to-left
 * modulo loop that yields natural front zero-padding.
 */
export class Base62Codec {
    private readonly _base62Alphabet: Alphabet;

    /**
     * Constructs a Base62Codec bound to the given Base62 alphabet.
     */
    public constructor(base62Alphabet: Alphabet) {
        this._base62Alphabet = base62Alphabet;
    }

    /**
     * Renders the given 32-bit value as a Base62 string of exactly the requested width, filling
     * positions from the least significant digit backward.
     */
    public encode(value: number, width: number): string {
        const radix: number = this._base62Alphabet.size();
        const characters: string[] = new Array<string>(width);

        let remaining: number = value >>> 0;

        for (let position: number = width - 1; position >= 0; position--) {
            const remainder: number = remaining % radix;
            characters[position] = this._base62Alphabet.characterAt(remainder);
            remaining = Math.floor(remaining / radix);
        }

        return characters.join("");
    }
}
