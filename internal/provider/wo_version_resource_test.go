package provider

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestFingerprint(t *testing.T) {
	salt := []byte("0123456789abcdef")
	h1, v1 := fingerprint("secret", salt)
	h2, v2 := fingerprint("secret", salt)
	if h1 != h2 || v1 != v2 {
		t.Fatal("fingerprint is not deterministic")
	}
	if h3, _ := fingerprint("secret", []byte("fedcba9876543210")); h3 == h1 {
		t.Fatal("salt has no effect")
	}
	if h4, _ := fingerprint("other", salt); h4 == h1 {
		t.Fatal("different values produced the same hash")
	}
	want, _ := strconv.ParseInt(h1[:8], 16, 64)
	if v1 != want {
		t.Fatalf("version %d != parseint(substr(hash, 0, 8), 16) %d", v1, want)
	}
}

const testAccWoVersionConfig = `
variable "secret" {
  type      = string
  ephemeral = true
}

resource "tools_wo_version" "test" {
  value_wo = var.secret
}
`

func TestAccWoVersion(t *testing.T) {
	res := "tools_wo_version.test"
	sameVersion := statecheck.CompareValue(compareSame)
	changedVersion := statecheck.CompareValue(compareDiffer)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_11_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:          testAccWoVersionConfig,
				ConfigVariables: config.Variables{"secret": config.StringVariable("first")},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(res, tfjsonpath.New("value_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(res, tfjsonpath.New("hash"), knownvalue.StringRegexp(hexRe)),
					sameVersion.AddStateValue(res, tfjsonpath.New("version")),
					changedVersion.AddStateValue(res, tfjsonpath.New("version")),
				},
			},
			{
				Config:          testAccWoVersionConfig,
				ConfigVariables: config.Variables{"secret": config.StringVariable("first")},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					sameVersion.AddStateValue(res, tfjsonpath.New("version")),
				},
			},
			{
				Config:          testAccWoVersionConfig,
				ConfigVariables: config.Variables{"secret": config.StringVariable("second")},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(res, tfjsonpath.New("version"), knownvalue.NotNull()),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					changedVersion.AddStateValue(res, tfjsonpath.New("version")),
				},
			},
		},
	})
}
