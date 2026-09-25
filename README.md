# terraform-provider-tools

Terraform provider with utility resources. Requires Terraform 1.11+.

## `tools_wo_version`

Turns a write-only/ephemeral value into a non-sensitive `version` that changes only when the value
changes — so `*_wo_version` arguments don't need to be bumped by hand.

```hcl
resource "tools_wo_version" "db_password" {
  value_wo = ephemeral.random_password.db.result
}

resource "aws_db_instance" "example" {
  password_wo         = ephemeral.random_password.db.result
  password_wo_version = tools_wo_version.db_password.version
}
```

A provider function can't do this: Terraform marks any function result derived from an ephemeral
value as ephemeral. A resource can read the write-only value from config at plan time and emit a
regular computed attribute.

## Development

```bash
go test ./... -v
```

Local override (`~/.terraformrc`):

```hcl
provider_installation {
  dev_overrides {
    "justereseau/tools" = "/path/to/go/bin"
  }
  direct {}
}
```

## Release

Releases are automated with [release-please](https://github.com/googleapis/release-please) from
[Conventional Commits](https://www.conventionalcommits.org/) (`fix:` → patch, `feat:` → minor, `feat!:` → major).

1. Merge PRs into `main` with Conventional Commit titles; release-please keeps a release PR up to date.
2. Merge the release PR **by hand**. The pipeline then runs tests, builds and signs with goreleaser, and
   creates the `vX.Y.Z` tag + GitHub Release, which the Terraform Registry picks up.

Pushing a `v*` tag manually also triggers a release. Requires the `GPG_PRIVATE_KEY` and `GPG_PASSPHRASE`
repo secrets and the public key registered on registry.terraform.io.
