#!/usr/bin/env node


import { Command } from "commander";
import { Cutter, SecurityContext, ValidationError } from "keycutter";


/**
 * Wires the Keycutter generation and validation pipelines behind a Commander-driven
 * command-line surface, exposing a token generator and a token verifier.
 */
export class Cli {
    private static readonly _PROGRAM_NAME: string = "keycutter";
    private static readonly _PROGRAM_VERSION: string = "1.0.0";

    private readonly _cutter: Cutter;
    private readonly _program: Command;

    public constructor() {
        this._cutter = new Cutter();
        this._program = new Command();

        this.configure();
    }

    private configure(): void {
        this._program
            .name(Cli._PROGRAM_NAME)
            .description("Command-line interface for the context-aware credential architecture, generating and verifying structured tokens.")
            .version(Cli._PROGRAM_VERSION);

        this._program
            .command("generate")
            .description("generate a structured token from the given semantic identifiers")
            .argument("<system>", "system identifier, restricted to lowercase ASCII letters and digits")
            .argument("<environment>", "environment identifier, restricted to lowercase ASCII letters and digits")
            .argument("<domain-purpose>", "domain purpose identifier, restricted to lowercase ASCII letters and digits")
            .option("-t, --timestamp", "embed a self-describing Base36 Unix-seconds creation timestamp", false)
            .action((
                systemIdentifier: string,
                environmentIdentifier: string,
                domainPurposeIdentifier: string,
                commandOptions: { timestamp: boolean },
            ): void => {
                this.generateToken(systemIdentifier, environmentIdentifier, domainPurposeIdentifier, commandOptions.timestamp);
            });

        this._program
            .command("verify")
            .description("verify a token and reveal its self-describing timestamp")
            .argument("<token>", "raw token to validate.")
            .action((rawToken: string): void => {
                this.verifyToken(rawToken);
            });
    }

    private generateToken(
        systemIdentifier: string,
        environmentIdentifier: string,
        domainPurposeIdentifier: string,
        embedTimestamp: boolean,
    ): void {
        try {
            const token: string = embedTimestamp
                ? this._cutter.generateWithTimestamp(systemIdentifier, environmentIdentifier, domainPurposeIdentifier)
                : this._cutter.generate(systemIdentifier, environmentIdentifier, domainPurposeIdentifier);

            process.stdout.write(`${token}\n`);
        } catch (error: unknown) {
            this.reportFailure(error);
        }
    }

    private verifyToken(rawToken: string): void {
        try {
            const context: SecurityContext = this._cutter.validate(rawToken);

            process.stdout.write("true\n");

            const createdAtUnixSeconds: bigint | null = context.createdAtUnixSeconds();

            if (createdAtUnixSeconds !== null) {
                process.stdout.write(`${createdAtUnixSeconds}\n`);
            }
        } catch (error: unknown) {
            this.reportFailure(error);
        }
    }

    private reportFailure(error: unknown): void {
        if (!(error instanceof ValidationError)) {
            throw error;
        }

        process.stderr.write(`${error.message}\n`);
        process.exitCode = 1;
    }

    public run(argumentVector: string[]): void {
        this._program.parse(argumentVector);
    }
}


new Cli().run(process.argv);
