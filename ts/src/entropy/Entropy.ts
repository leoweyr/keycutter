/**
 * The fixed-width high-intensity random segment of a token.
 */
export class Entropy {
    /**
     * The fixed number of Base62 characters in the high-intensity entropy segment, sized to deliver
     * roughly 143 bits of randomness.
     */
    public static readonly LENGTH: number = 24;

    private readonly _value: string;

    public constructor(value: string) {
        this._value = value;
    }

    /** @returns The raw entropy characters. */
    public getValue(): string {
        return this._value;
    }

    /** @returns The number of characters in the entropy segment. */
    public length(): number {
        return this._value.length;
    }

    /** @returns The raw entropy characters. */
    public toString(): string {
        return this._value;
    }
}
