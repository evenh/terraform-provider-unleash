package unleash

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/evenh/terraform-provider-unleash/unleash/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.RunWithUnleash(func(unleashPort int) int {
		_ = os.Setenv(UNLEASH_API_ENDPOINT, fmt.Sprintf("http://localhost:%d/api", unleashPort))
		_ = os.Setenv(UNLEASH_AUTH_EMAIL, "acceptance-test@unleash.provider.tf")

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
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-"+id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "It works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-"+id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "It works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "true"),
				),
			},
			{
				Config: testAccUnleashFeatureToggle_update(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-"+id),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "Update works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "false"),
				),
			},
		},
	})
}

func testAccUnleashFeatureToggle_basic(id string) string {
	return fmt.Sprintf(`
resource "unleash_feature_toggle" "basic" {
  name        = "acctest-%s"
  description = "It works"
  enabled     = true
}
`, id)
}

func testAccUnleashFeatureToggle_update(id string) string {
	return fmt.Sprintf(`
resource "unleash_feature_toggle" "basic" {
  name        = "acctest-%s"
  description = "Update works"
  enabled     = false
}
`, id)
}
