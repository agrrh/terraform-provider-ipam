// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAllocationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAllocationResourceConfig("foo", 24),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_allocation.foo", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ipam_allocation.foo", "id", "10.0.0.0/24"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ipam_allocation.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"cidr", "pool_id", "size"},
			},
			// Update and Read testing
			{
				Config: testAccAllocationResourceConfig("bar", 24),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ipam_allocation.bar", "cidr", "10.0.0.0/24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAllocationResourceConfig(name string, configurableAttribute int) string {
	return fmt.Sprintf(`
resource "ipam_pool" "test" {
  cidr = "10.0.0.0/16"
}

resource "ipam_allocation" %[1]q {
  pool_id = ipam_pool.test.id
	size    = %[2]d
}
`, name, configurableAttribute)
}
