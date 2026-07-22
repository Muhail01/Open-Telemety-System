#!/usr/bin/env sh
set -eu

fail=0

scan() {
  label="$1"
  pattern="$2"
  if git grep -nEI "$pattern" -- ':!scripts/audit_public_repo.sh' > /tmp/public-audit-matches.txt 2>/dev/null; then
    echo "ERROR: public repository audit matched $label"
    cat /tmp/public-audit-matches.txt
    fail=1
  fi
}

# Common credential material. The audit script itself is excluded so the
# detection expressions do not match their own source.
scan "private key material" '-----BEGIN (RSA|OPENSSH|EC|DSA) PRIVATE KEY-----'
scan "AWS access key" '(AKIA|ASIA)[0-9A-Z]{16}'
scan "Google API key" 'AIza[0-9A-Za-z_-]{30,}'
scan "Google OAuth client secret" 'GOCSPX-[0-9A-Za-z_-]{20,}'
scan "OpenAI-style secret" 'sk-[A-Za-z0-9_-]{20,}'
scan "bearer token" 'Bearer[[:space:]]+[A-Za-z0-9._-]{20,}'
scan "VLESS connection URI" 'vless://'

# Known private GILLZY deployment identifiers that must not leak into the
# standalone public implementation.
scan "private GILLZY production domain" 'gillzy-store\.com'
scan "private GILLZY server path" '/opt/gillzy'
scan "private GILLZY LAN address" '192\.168\.1\.30'
scan "private monolith import path" 'marketplace/internal/'

rm -f /tmp/public-audit-matches.txt

if [ "$fail" -ne 0 ]; then
  exit 1
fi

echo "Public repository secret/boundary audit passed."
