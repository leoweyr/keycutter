/**
 * Utf8 converts text into its UTF-8 byte stream, mirroring Go's []byte(string) conversion so that
 * byte-wise operations stay identical across language ports.
 */
export class Utf8 {
    private static readonly _encoder: TextEncoder = new TextEncoder();

    /**
     * Encodes the given text as its UTF-8 byte stream.
     */
    public static encode(text: string): Uint8Array {
        return Utf8._encoder.encode(text);
    }
}
