package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePrefixedUniqueId(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrefixedUniqueId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"internals_prefixed_unique_id.foo", "id", regexp.MustCompile("^bar\\d{18}[0-9a-f]{8}$")),
				),
			},
		},
	})
}

const testAccResourcePrefixedUniqueId = `
resource "internals_prefixed_unique_id" "foo" {
  prefix = "bar"
}
`
