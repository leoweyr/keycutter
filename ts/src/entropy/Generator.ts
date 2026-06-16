import { Entropy } from "./Entropy";


/**
 * Generator produces a fresh high-intensity entropy segment on demand.
 */
export interface Generator {
    generate(): Entropy;
}
