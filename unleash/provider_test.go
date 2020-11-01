package unleash

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This map is most commonly constructed once in a common init() method of the Provider’s main test file,
// and includes an object of the current Provider type. https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html
var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"unleash": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}
