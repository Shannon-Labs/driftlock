# Driftlock Commercial License

This license governs use of the Driftlock API binaries and containers ("Commercial Software"). The CBAD Rust core remains licensed under Apache 2.0.

## Grant
- Shannon Labs grants Customer a non-transferable, non-exclusive right to run the Commercial Software for internal evaluation or production, subject to an active license key.
- Redistribution is prohibited unless explicitly approved in writing.

## Activation
- Each deployment must provide `DRIFTLOCK_LICENSE_KEY` at startup. The key encodes tenant, tier, and expiry.
- Keys are signed by Shannon Labs. The server validates the signature and refuses to start if the key is missing, invalid, or expired.
- Evaluation keys (`EVAL-*`) expire automatically; production keys must be renewed before the expiry timestamp embedded in the key.

## Dev Mode (Local Demos Only)
- Setting `DRIFTLOCK_DEV_MODE=true` bypasses the license check to simplify local onboarding.
- `/healthz` clearly reports `license.status="dev_mode"` while the bypass is active; any production or pilot deployment must use a signed key instead.
- Shannon Labs reserves the right to disable dev mode in environments that expose it publicly.

## Restrictions
- Do not attempt to circumvent license enforcement or modify signatures.
- Do not provide the Commercial Software as a managed service to third parties without a reseller agreement.
- Do not reverse engineer or decompile the Commercial Software except where permitted by law.

## Compliance
- Customer must maintain records of deployed instances and make them available to Shannon Labs upon request for compliance verification.
- Shannon Labs may revoke or refuse to renew licenses if the Customer violates these terms.

## Warranty & Liability
- Commercial Software is provided "as is" without warranties. Shannon Labs is not liable for indirect or consequential damages.

## Contact
For licensing questions or additional terms (OEM, reseller, or white-label agreements), contact legal@shannonlabs.com.
