# Tokenforge

One spec, every language — context-aware credential architecture for generating and verifying structured tokens with byte-identical layout and CRC32 tail checksums.

```
[SystemIdentifier]_[EnvironmentIdentifier]_[DomainPurposeIdentifier]_(Timestamp_)[Entropy][Checksum]
```

| Segment                | Content                                                                                                                                                                                                                                                                                                                                        | Width          |
|------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------|
| Prefix                 | System, environment, domain purpose identifiers, optional timestamp, delimited by `_`<br />Restricts every semantic component to lowercase ASCII letters and decimal digits only. Underscores, uppercase letters, and any multi-byte non-ASCII characters are strictly forbidden.<br />`BASE36_CHARSET = 0123456789abcdefghijklmnopqrstuvwxyz` | Variable (≥ 6) |
| High-intensity Entropy | 24 Base62 characters<br />Strictly the 62-character GMP dictionary, ordered by ASCII code in ascending order. Third-party Base62 dialects that embed internal permutations are prohibited<br />`BASE62_CHARSET = 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`                                                               | 24             |
| Tail Checksum          | 6 Base62 characters encoding CRC32-IEEE of Prefix + Entropy                                                                                                                                                                                                                                                                                    | 6              |

**Physical Seamless Fusion** — No separator is permitted between the high-intensity entropy and the tail checksum. The full-stack slice cursor operates as a hard asymmetric dead-lock: the last 6 characters are always the checksum, everything before them is always the base string.

**Absolute Unambiguous Prefix** — Each semantic component of the prefix is subject to strict character-set constraints that fundamentally eliminate cross-platform parsing ambiguity and the hash inconsistency introduced by Unicode normalization (NFC/NFD).

**Unbiased Uniform Sampling** — Direct modulo on a low-bit-width integer is prohibited across all platforms. Every random index must be drawn from an OS-level cryptographically secure random number generator (CSPRNG) via rejection sampling, completely eliminating modulo bias.

**Structured Threshold Assertion** — Perimeter length guards discard hard-coded magic numbers. Minimum valid lengths are derived dynamically from the credential's own topology formula.

**Self-Describing Key Versioning** *(optional)* — An optional monotonic Unix-seconds timestamp embeds each credential's own creation instant, imposing a natural total order over key versions so that rotation, age-based expiry, and which-key-is-newer all resolve from the tokens alone, with no central registry lookup.

## 🚀 Quick Start

### Go

```bash
go get go.leoweyr.com/tokenforge/go
```

## 🏗️ Generation Pipeline

### 1. Prefix Construction

Concatenate plaintext semantic segments for the given use case. Prefix length is variable and determined by system design, but must conform to a fixed topology.

*Example*: `odc_prod_msk`

Render the current Unix epoch seconds in Base36 and append it as a fourth prefix component after the domain purpose identifier. Base36 keeps every digit inside `BASE36_CHARSET`, and it is variable-width — 6 characters through late 2038, then 7.

> [!WARNING]
>
> Embed it only for keys inside a trust boundary, since the creation instant is readable by any holder.

$$
\text{Timestamp} = \text{Base36}(\text{UnixSeconds})
$$

*Example*: `odc_prod_msk_tgny7c`

### 2. Unbiased Entropy Generation

Draw 24 characters from `BASE62_CHARSET` using an OS-level CSPRNG with a rejection-sampling loop to ensure each index is drawn uniformly from `[0, 61]`. Direct modulo on a raw random integer is forbidden — it introduces a statistical skew that taints the distribution. The output is a strictly fixed 24-character string.

*Example*: `7xT2zP9qL4wK1mN8vV5cB3nA`

### 3. Base String Concatenation

Concatenate the prefix from step 1 and the high-intensity entropy string from step 2 directly.

$$
\text{BaseString} = \text{Prefix} + \text{HighIntensityEntropy}
$$

*Example (no timestamp)*: `odc_prod_msk_7xT2zP9qL4wK1mN8vV5cB3nA`

*Example (with timestamp)*: `odc_prod_msk_tgny7c_7xT2zP9qL4wK1mN8vV5cB3nA`

### 4. Mathematical CRC32 Mapping

Compute the CRC32-IEEE (reflected form) checksum of `BaseString` encoded as a UTF-8 byte stream, yielding a 32-bit unsigned integer $\text{Value}$. Convert $\text{Value}$ to a 6-character Base62 string $\text{Checksum}$ using the following right-to-left modulo loop:

$$
\begin{array}{l}
\text{for } i = 5 \rightarrow 0: \\
\quad \text{Remainder} = \text{Value} \bmod 62 \\
\quad \text{Checksum}[i] = \text{BASE62_CHARSET}[\text{Remainder}] \\
\quad \text{Value} = \left\lfloor \dfrac{\text{Value}}{62} \right\rfloor
\end{array}
$$

The loop fills from the last position backward. Any value that does not require all 6 digits is naturally zero-padded at the front — no explicit padding logic is needed.

### 5. Final Assembly

Append the 6-character checksum directly to the end of the base string with no separator:

$$
\text{Token} = \text{BaseString} + \text{TailChecksum}
$$

*Example*: `odc_prod_msk_7xT2zP9qL4wK1mN8vV5cB3nA4VHrHM`

## 🛡️ Validation Pipeline

### 1. Structural Guard & Asymmetric Slice

At every network edge gateway or application entry point, perform hard physical boundary assertions before any business logic runs.

Length check:

$$
\text{MinLength} = \text{Len}(\text{Prefix}) + 24 + 6
$$

If the validator holds a known expected `Prefix`, the token length must equal exactly $\text{len}(\text{Prefix}) + 30$. Without a known prefix, the absolute minimum token length is 36 characters (each of the three prefix components must be at least 1 character, so the shortest valid prefix is 6 characters). Any token shorter than the derived threshold is discarded immediately.

Atomic slice, ignoring internal underscore structure, using the fixed-width checksum cursor:

$$
\begin{matrix}
\text{BaseString} = \text{Token}[0:\text{len}(\text{Token}) - 6] \\
\text{ProvidedTailCheckSum} = \text{Token}[\text{len}(\text{Token}) - 6:\text{len}(\text{Token})]
\end{matrix}
$$

### 2. Idempotent Verification

Re-execute the generation pipeline step 4 mapping locally on the extracted `BaseString` to obtain `ExpectedTailChecksum`. If `ExpectedChecksum ≠ ProvidedChecksum`, the token is rejected immediately as corrupted or truncated — fail-fast, no further processing.

### 3. Context Reification

Only after passing step 2 may the system split the prefix portion of `BaseString` on `_`. Because `BASE36_CHARSET` forbids underscores within any component, the split result is uniquely deterministic across every platform and encoding. Extract the system, environment, and domain purpose identifiers and inject them as a typed security context object into downstream operations.

The prefix component count determines whether an optional timestamp is present: the validator is never told in advance, three components carry no timestamp, four mark the trailing one as the Base36 timestamp.

## ⚖️ Why 24 Characters

On the cryptographic side, 24 Base62 characters carry:

$$
62^{24} \approx 2^{143} \text{ bits}
$$

The practical security floor for API tokens in cloud-native architecture is 128 bits. 24 characters delivers 143 bits — clearing the threshold with margin to spare.

Some vendors push further: GitHub's personal access tokens use 30 entropy characters, reaching ~178 bits. From a pure cryptographic standpoint, that number is unimpeachable. From an engineering leverage standpoint, it is unnecessary — beyond ~140 bits, the marginal security return per additional character converges to zero. The cost, however, is real. Every extra character widens network payloads, inflates database index pages, and — on mobile — turns a token into a string too long to select cleanly with a long-press. 

Tokenforge locks high-intensity entropy at 24 characters: the precise point where cryptographic surplus meets transmission efficiency and human ergonomics, with nothing wasted on either side.
