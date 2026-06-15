import { Alphabet } from "../encoding/Alphabet";
import { Entropy } from "./Entropy";
import { Generator } from "./Generator";


/**
 * SecureGenerator draws entropy from an operating-system CSPRNG and applies rejection sampling to
 * guarantee an unbiased uniform symbol distribution.
 */
export class SecureGenerator implements Generator {
    private readonly _base62Alphabet: Alphabet;

    public constructor(base62Alphabet: Alphabet) {
        this._base62Alphabet = base62Alphabet;
    }

    /**
     * Produces a fixed-width entropy segment by drawing unbiased indices from the bound alphabet. A
     * 32-byte bulk read reduces CSPRNG syscall overhead to O(1).
     */
    public generate(): Entropy {
        const characters: string[] = new Array<string>(Entropy.LENGTH);
        const size: number = this._base62Alphabet.size();
        const ceiling: number = 256 - (256 % size);
        const buffer: Uint8Array = new Uint8Array(32);

        let bufferPosition: number = buffer.length;
        let position: number = 0;

        while (position < Entropy.LENGTH) {
            if (bufferPosition >= buffer.length) {
                crypto.getRandomValues(buffer);
                bufferPosition = 0;
            }

            const candidate: number = buffer[bufferPosition]!;
            bufferPosition++;

            if (candidate < ceiling) {
                characters[position] = this._base62Alphabet.characterAt(candidate % size);
                position++;
            }
        }

        return new Entropy(characters.join(""));
    }
}
