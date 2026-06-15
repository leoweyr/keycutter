import { Checksum } from "./Checksum";


/**
 * Calculator maps a token base string to its fixed-width Base62 tail checksum.
 */
export interface Calculator {
    calculate(baseString: string): Checksum;
}
