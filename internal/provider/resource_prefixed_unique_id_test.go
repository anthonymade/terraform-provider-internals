package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestAccResourcePrefixedUniqueIdImport(t *testing.T) {
	importId := "prefix,20220107181856110700000001"
	expectedId := "prefix20220107181856110700000001"
	expectedPrefix := "prefix"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:           testAccResourcePrefixedUniqueId,
				ImportState:      true,
				ResourceName:     "internals_prefixed_unique_id.foo",
				ImportStateId:    importId,
				ImportStateCheck: checkImportFunc(expectedId, expectedPrefix),
			},
		},
	})
}

func TestAccResourcePrefixedUniqueIdImportNoPrefix(t *testing.T) {
	importId := "20220107181856110700000001"
	expectedId := "20220107181856110700000001"
	expectedPrefix := ""
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:           testAccResourcePrefixedUniqueId,
				ImportState:      true,
				ResourceName:     "internals_prefixed_unique_id.foo",
				ImportStateId:    importId,
				ImportStateCheck: checkImportFunc(expectedId, expectedPrefix),
			},
		},
	})
}

func TestAccResourcePrefixedUniqueIdImportMultiSeparators(t *testing.T) {
	importId := "why,why,not,20220107181856110700000001"
	expectedId := "why,why,not20220107181856110700000001"
	expectedPrefix := "why,why,not"
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:           testAccResourcePrefixedUniqueId,
				ImportState:      true,
				ResourceName:     "internals_prefixed_unique_id.foo",
				ImportStateId:    importId,
				ImportStateCheck: checkImportFunc(expectedId, expectedPrefix),
			},
		},
	})
}

const testAccResourcePrefixedUniqueId = `
resource "internals_prefixed_unique_id" "foo" {
  prefix = "bar"
}
`

func checkImportFunc(expectedId string, expectedPrefix string) func(states []*terraform.InstanceState) error {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("want 1 state, got: %+v", states)
		}
		state := states[0]
		if state.ID != expectedId {
			return fmt.Errorf("ID want %v, got: %v", expectedId, state.ID)
		}
		if state.Attributes["prefix"] != expectedPrefix {
			return fmt.Errorf("prefix want %v, got: %v", expectedPrefix, state.Attributes["prefix"])
		}
		return nil
	}
}
