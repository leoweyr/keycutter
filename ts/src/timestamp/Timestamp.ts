/**
 * The optional creation marker embedded as the fourth prefix component. It pairs a Unix epoch
 * second count with its Base36 rendering, whose digits remain inside the permitted prefix character
 * set.
 */
export class Timestamp {
    private readonly _seconds: bigint;
    private readonly _encoded: string;

    public constructor(seconds: bigint, encoded: string) {
        this._seconds = seconds;
        this._encoded = encoded;
    }

    /** @returns The embedded Unix epoch second count. */
    public getSeconds(): bigint {
        return this._seconds;
    }

    /** @returns The Base36 rendering placed inside the token prefix. */
    public getValue(): string {
        return this._encoded;
    }

    /** @returns The Base36 rendering placed inside the token prefix. */
    public toString(): string {
        return this._encoded;
    }
}
