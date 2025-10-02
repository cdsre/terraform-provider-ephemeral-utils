// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRevealerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRevealerResourceConfig("hello"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"ephemeral-utils_revealer.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(".+")),
					),
					statecheck.ExpectKnownValue(
						"ephemeral-utils_revealer.test",
						tfjsonpath.New("data"),
						knownvalue.StringExact("hello"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:            "ephemeral-utils_revealer.test",
				ImportState:             true,
				ImportStateId:           "hello",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
			},
			// Update and Read testing
			{
				Config: testAccRevealerResourceConfig("world"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"ephemeral-utils_revealer.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(".+")),
					),
					statecheck.ExpectKnownValue(
						"ephemeral-utils_revealer.test",
						tfjsonpath.New("data"),
						knownvalue.StringExact("world"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRevealerResourceConfig(val string) string {
	return fmt.Sprintf(`
resource "ephemeral-utils_revealer" "test" {
  data_wo = %[1]q
}
`, val)
}
