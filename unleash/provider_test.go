package unleash

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This map is most commonly constructed once in a common init() method of the Provider’s main test file,
// and includes an object of the current Provider type. https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html
var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"unleash": testAccProvider,
	}
}
