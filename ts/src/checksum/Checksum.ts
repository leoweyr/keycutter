/**
 * The fixed-width Base62 tail segment encoding the CRC32 integrity value of a token base string.
 */
export class Checksum {
    /**
     * The fixed number of Base62 characters in the tail checksum segment, wide enough to encode any
     * 32-bit CRC value without truncation.
     */
    public static readonly LENGTH: number = 6;

    private readonly _value: string;

    public constructor(value: string) {
        this._value = value;
    }

    /** @returns The raw checksum characters. */
    public getValue(): string {
        return this._value;
    }

    /**
     * @returns Reports whether this checksum carries the same characters as the other checksum.
     */
    public equals(other: Checksum): boolean {
        return this._value === other._value;
    }

    /** @returns The raw checksum characters. */
    public toString(): string {
        return this._value;
    }
}
