/**
 * Crc32 computes the CRC32-IEEE checksum, reproducing Go's hash/crc32.ChecksumIEEE bit for bit so
 * the integrity value stays identical across every Keycutter language port.
 */
export class Crc32 {
    /**
     * The reversed CRC32-IEEE generator polynomial (0xEDB88320), the same constant Go's hash/crc32
     * uses for the IEEE table. It matches the IEEE 802.3 definition.
     */
    private static readonly _POLYNOMIAL: number = 0xedb88320;

    /**
     * The 256-entry lookup table derived from the IEEE polynomial, built once at class load to
     * mirror Go's package-level crc32.IEEETable.
     */
    private static readonly _TABLE: Uint32Array = Crc32.buildTable();

    /**
     * Precomputes the byte-wise CRC32 lookup table for the IEEE polynomial.
     */
    private static buildTable(): Uint32Array {
        const table: Uint32Array = new Uint32Array(256);

        for (let index: number = 0; index < 256; index++) {
            let remainder: number = index;

            for (let bit: number = 0; bit < 8; bit++) {
                if ((remainder & 1) !== 0) {
                    remainder = (remainder >>> 1) ^ Crc32._POLYNOMIAL;
                } else {
                    remainder = remainder >>> 1;
                }
            }

            table[index] = remainder >>> 0;
        }

        return table;
    }

    /**
     * Computes the CRC32-IEEE value of the given byte stream and returns it as an unsigned 32-bit
     * integer.
     */
    public static checksumIeee(bytes: Uint8Array): number {
        let crc: number = 0xffffffff;

        for (let position: number = 0; position < bytes.length; position++) {
            const tableIndex: number = (crc ^ bytes[position]!) & 0xff;
            crc = (crc >>> 8) ^ Crc32._TABLE[tableIndex]!;
        }

        return (crc ^ 0xffffffff) >>> 0;
    }
}
