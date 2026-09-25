---
page_title: "tools Provider"
description: |-
  Utility resources for Terraform workflows.
---

# tools Provider

Utility resources for Terraform workflows. Requires Terraform 1.11+.

Resource names start with `jr_`, so declare the provider with the local name `jr`.

## Example Usage

```terraform
terraform {
  required_providers {
    jr = {
      source = "justereseau/tools"
    }
  }
}
```
