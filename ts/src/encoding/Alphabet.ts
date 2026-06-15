import { Utf8 } from "./Utf8";


/**
 * An immutable ordered symbol set that maps between character values and their positional indices.
 */
export class Alphabet {
    /**
     * The canonical 62-symbol GMP dictionary ordered by ascending ASCII code, used for
     * high-intensity entropy and tail checksum encoding.
     */
    public static readonly BASE62_CHARACTERS: string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

    /**
     * The canonical 36-symbol dictionary ordered by ascending ASCII code, spanning the decimal
     * digits and lowercase ASCII letters.
     */
    public static readonly BASE36_CHARACTERS: string = "0123456789abcdefghijklmnopqrstuvwxyz";

    private readonly _characters: string;
    private readonly _indexByCharacter: Map<number, number>;

    /**
     * Constructs an Alphabet from the given ordered character set and indexes every symbol for
     * constant-time membership and lookup queries.
     */
    public constructor(characters: string) {
        const indexByCharacter: Map<number, number> = new Map<number, number>();

        for (let position: number = 0; position < characters.length; position++) {
            indexByCharacter.set(characters.charCodeAt(position), position);
        }

        this._characters = characters;
        this._indexByCharacter = indexByCharacter;
    }

    /** @returns The number of symbols contained in the alphabet. */
    public size(): number {
        return this._characters.length;
    }

    /**
     * @returns The symbol located at the given positional index.
     */
    public characterAt(index: number): string {
        return this._characters.charAt(index);
    }

    /**
     * @returns Reports whether the given byte value belongs to the alphabet.
     */
    public contains(byteValue: number): boolean {
        return this._indexByCharacter.has(byteValue);
    }

    /**
     * @returns The positional index of the given byte value, or undefined when that byte value does
     * not belong to the alphabet.
     */
    public indexOf(byteValue: number): number | undefined {
        return this._indexByCharacter.get(byteValue);
    }

    /**
     * @returns Reports whether every character in the text belongs to the alphabet.
     */
    public permits(text: string): boolean {
        const bytes: Uint8Array = Utf8.encode(text);

        for (let position: number = 0; position < bytes.length; position++) {
            if (!this.contains(bytes[position]!)) {
                return false;
            }
        }

        return true;
    }
}
