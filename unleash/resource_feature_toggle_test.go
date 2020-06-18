package unleash

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/evenh/terraform-provider-unleash/unleash/internal/test"
)

func TestMain(m *testing.M) {
	test.RunWithUnleash(func() int {
		return m.Run()
	})
}

func TestOne(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnleashFeatureToggle_createAToggle(),
			},
		},
	})
}

func testAccUnleashFeatureToggle_createAToggle() string {
	return fmt.Sprintf(`
provider "unleash" {
  auth {
    unsecure {
       username = "acceptance@test.com"
    }
  }
}
resource "unleash_feature_toggle" "test" {
  name     = "acctest-%d"
  enabled = true
}
`, 42)
}
