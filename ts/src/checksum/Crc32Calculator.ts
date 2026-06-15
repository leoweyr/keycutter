import { Utf8 } from "../encoding/Utf8";
import { Base62Codec } from "../encoding/Base62Codec";
import { Calculator } from "./Calculator";
import { Checksum } from "./Checksum";
import { Crc32 } from "./Crc32";


/**
 * Crc32Calculator derives the tail checksum from the CRC32-IEEE value of a base string encoded as
 * Base62.
 */
export class Crc32Calculator implements Calculator {
    private readonly _encoder: Base62Codec;

    /**
     * Constructs a Crc32Calculator backed by the given Base62 codec.
     */
    public constructor(encoder: Base62Codec) {
        this._encoder = encoder;
    }

    /**
     * Computes the CRC32-IEEE value of the base string interpreted as a UTF-8 byte stream and
     * encodes it as a fixed-width Base62 checksum.
     */
    public calculate(baseString: string): Checksum {
        const value: number = Crc32.checksumIeee(Utf8.encode(baseString));
        const encoded: string = this._encoder.encode(value, Checksum.LENGTH);

        return new Checksum(encoded);
    }
}
