import { Timestamp } from "../timestamp/Timestamp";


/**
 * Bundles the three semantic identifiers and the optional timestamp extracted from a
 * token prefix.
 */
export interface ReifiedContext {
    systemIdentifier: string;
    environmentIdentifier: string;
    domainPurposeIdentifier: string;
    timestamp: Timestamp | null;
}
