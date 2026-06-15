import { Clock } from "./Clock";


/**
 * SystemClock reads the current instant from the operating-system wall clock.
 */
export class SystemClock implements Clock {
    private static _instance: SystemClock;

    public static getInstance(): SystemClock {
        if (!SystemClock._instance) {
            SystemClock._instance = new SystemClock();
        }

        return SystemClock._instance;
    }

    private constructor() {}

    /**
     * @returns The current Unix epoch second count from the system wall clock.
     */
    public nowUnixSeconds(): bigint {
        return BigInt(Math.floor(Date.now() / 1000));
    }
}
