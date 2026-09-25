---
page_title: "jrtools_wo_version Resource - jrtools"
description: |-
  Derives a stable, non-sensitive version from a write-only value, for use as a *_wo_version argument.
---

# jrtools_wo_version (Resource)

Derives a stable, non-sensitive version from a write-only value, for use as a `*_wo_version` argument.
The version changes only when the value changes, so the write-only value is re-sent only when needed.

The value is fingerprinted with Argon2id and a random per-resource salt, so the hash in state cannot be
cheaply brute-forced back to the secret.

On create, `hash` and `version` are known only after apply. On later plans they are computed from the
stored salt, so a change to the value shows up in the same plan.

## Example Usage

```terraform
ephemeral "random_password" "db" {
  length = 32
}

resource "jrtools_wo_version" "db_password" {
  value_wo = ephemeral.random_password.db.result
}

resource "aws_db_instance" "example" {
  # ...
  storage_encrypted   = true
  password_wo         = ephemeral.random_password.db.result
  password_wo_version = jrtools_wo_version.db_password.version
}
```

## Schema

### Required

- `value_wo` (String, Write-only) Value to fingerprint. Accepts ephemeral values; never stored in plan or state.

### Read-Only

- `hash` (String) Hex-encoded Argon2id hash of `value_wo`.
- `salt` (String) Random base64 salt generated on create.
- `version` (Number) First 8 hex characters of `hash` as an integer (0 to 4294967295).
