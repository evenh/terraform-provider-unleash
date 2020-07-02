package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/evenh/terraform-provider-unleash/unleash/internal/test"
)

var port = 4242

func TestMain(m *testing.M) {
	test.RunWithUnleash(func(unleashPort int) int {
		port = unleashPort
		return m.Run()
	})
}

func TestOne(t *testing.T) {
	resourceName := "unleash_feature_toggle.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashFeatureToggle_createAToggle(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, NAME, "acctest-42"),
					resource.TestCheckResourceAttr(resourceName, DESCRIPTION, "It works"),
					resource.TestCheckResourceAttr(resourceName, ENABLED, "true"),
				),
			},
		},
	})
}

func testAccUnleashFeatureToggle_createAToggle() string {
	return fmt.Sprintf(`
provider "unleash" {
  api_endpoint = "http://localhost:%d/api"

  auth {
    unsecure {
       email = "acceptance-test@unleash.provider.tf"
    }
  }
}
resource "unleash_feature_toggle" "test" {
  name        = "acctest-%d"
  description = "It works"
  enabled     = true
}
`, port, 42)
}
