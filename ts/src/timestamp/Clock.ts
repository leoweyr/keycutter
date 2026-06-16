/**
 * Clock supplies the current wall-clock instant used to stamp a token at generation time.
 */
export interface Clock {
    nowUnixSeconds(): bigint;
}
