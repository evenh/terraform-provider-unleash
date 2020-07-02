package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/evenh/terraform-provider-unleash/unleash/internal/test"
)

var port int

func TestMain(m *testing.M) {
	test.RunWithUnleash(func(unleashPort int) int {
		port = unleashPort
		return m.Run()
	})
}

func TestAccUnleashFeatureToggle_basic(t *testing.T) {
	id := acctest.RandString(6)
	resourceName := "unleash_feature_toggle.basic"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashFeatureToggle_basic(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-" + id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "It works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "true"),
				),
			},
		},
	})
}

func TestAccUnleashFeatureToggle_update(t *testing.T) {
	id := acctest.RandString(6)
	resourceName := "unleash_feature_toggle.basic"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashFeatureToggle_basic(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-" + id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "It works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "true"),
				),
			},
			{
				Config: testAccUnleashFeatureToggle_update(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-" + id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "Update works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "false"),
				),
			},
		},
	})
}

func testAccUnleashFeatureToggle_basic(id string) string {
	return fmt.Sprintf(`
%s

resource "unleash_feature_toggle" "basic" {
  name        = "acctest-%s"
  description = "It works"
  enabled     = true
}
`, providerBlock(), id)
}

func testAccUnleashFeatureToggle_update(id string) string {
	return fmt.Sprintf(`
%s

resource "unleash_feature_toggle" "basic" {
  name        = "acctest-%s"
  description = "Update works"
  enabled     = false
}
`, providerBlock(), id)
}

func providerBlock() string {
	return fmt.Sprintf(`
provider "unleash" {
  api_endpoint = "http://localhost:%d/api"

  auth {
    unsecure {
       email = "acceptance-test@unleash.provider.tf"
    }
  }
}
`, port)
}
