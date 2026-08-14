# Release verification

Treecat releases publish checksums, an SPDX JSON SBOM, and a GitHub artifact attestation for the checksums manifest.

## Verify a downloaded artifact

Download the artifact and `checksums.txt` from the same GitHub release, then run:

```bash
sha256sum --check --ignore-missing checksums.txt
```

On macOS, use:

```bash
shasum -a 256 -c checksums.txt
```

The command must report the downloaded artifact as `OK` before installation.

## Verify build provenance

GitHub's attestation is generated in the release workflow with OIDC-backed signing. With the GitHub CLI installed, verify the checksums manifest against this repository:

```bash
gh attestation verify checksums.txt --repo carterlasalle/treecat
```

A valid result proves that the manifest was produced by a GitHub Actions workflow in `carterlasalle/treecat`; it does not replace checking the manifest against the artifact itself.

## Inspect the SBOM

Each release also uploads `treecat.spdx.json`. It describes the release output in SPDX JSON form and can be inspected directly or with SPDX-aware tooling.

## Release pipeline controls

The release workflow:

1. verifies and tidies the Go module graph without accepting changes;
2. runs the race-enabled test suite and `govulncheck`;
3. builds all configured GoReleaser artifacts with pinned tooling;
4. checks every generated checksum;
5. smoke-tests the Linux binary, Windows ZIP contents, and Debian install/uninstall path;
6. generates and uploads the SBOM; and
7. publishes a provenance attestation for `checksums.txt`.

Third-party GitHub Actions and release tools are pinned to immutable revisions or explicit versions so releases do not silently change underneath an existing tag.
