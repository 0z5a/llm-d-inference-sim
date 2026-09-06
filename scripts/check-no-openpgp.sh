#!/usr/bin/env bash
# GO-2026-5932 affects the unmaintained OpenPGP package, not all of x/crypto.
set -euo pipefail

packages=$(go list -mod=readonly -deps -test ./...)
if [ -z "$packages" ]; then
  echo "ERROR: no Go packages found." >&2
  exit 1
fi

if matches=$(printf '%s\n' "$packages" | grep -E '^golang\.org/x/crypto/openpgp(/|$)'); then
  echo "ERROR: deprecated OpenPGP dependency detected (GO-2026-5932)." >&2
  printf '%s\n' "$matches" >&2
  exit 1
fi

echo "No deprecated OpenPGP packages found in dependencies (including tests)."
