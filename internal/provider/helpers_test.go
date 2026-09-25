package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

var (
	hexRe         = regexp.MustCompile(`^[0-9a-f]{64}$`)
	compareSame   = compare.ValuesSame()
	compareDiffer = compare.ValuesDiffer()
)
