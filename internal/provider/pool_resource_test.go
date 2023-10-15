// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPoolResourceConfig("foo", "10.0.0.0/16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_pool.foo", "cidr", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("ipam_pool.foo", "id", "10.0.0.0/16"),
				),
			},
			{
				Config: testAccPoolResourceConfig("bar", "10.5.0.0/16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_pool.bar", "cidr", "10.5.0.0/16"),
					resource.TestCheckResourceAttr("ipam_pool.bar", "id", "10.5.0.0/16"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ipam_pool.bar",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"cidr"},
			},
			// Update and Read testing
			{
				Config: testAccPoolResourceConfig("foo", "10.10.0.0/16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_pool.foo", "cidr", "10.10.0.0/16"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPoolResourceConfig(name, configurableAttribute string) string {
	return fmt.Sprintf(`
resource "ipam_pool" %[1]q {
  cidr = %[2]q
}
`, name, configurableAttribute)
}
