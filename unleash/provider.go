package unleash

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/meta"

	"github.com/evenh/terraform-provider-unleash/unleash/api"
	"github.com/evenh/terraform-provider-unleash/version"
)

// Environment Variables
const (
	UNLEASH_API_ENDPOINT = "UNLEASH_API_ENDPOINT"
)

// Provider keys
const (
	api_endpoint           = "api_endpoint"
	auth                   = "auth"
	auth_unsecure          = "unsecure"
	auth_unsecure_username = "username"
	auth_unsecure_email    = "email"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			api_endpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(UNLEASH_API_ENDPOINT, "http://localhost:4242/api"),
				Description: "The Unleash API endpoint, e.g. http://localhost:4242/api",
			},
			auth: {
				Type:     schema.TypeSet,
				Optional: false,
				MinItems: 1,
				MaxItems: 1,
				Description: "Authentication mechanism to use for communicating with the Unleash API",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						auth_unsecure: {
							Type:     schema.TypeSet,
							Required: false,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									auth_unsecure_email: {
										Type: schema.TypeString,
									},
									auth_unsecure_username: {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringIsNotWhiteSpace,
									},
								},
							},
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"unleash_feature_toggle": resourceFeatureToggle(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseUrl := d.Get(api_endpoint).(string)
	userAgent := fmt.Sprintf("terraform-provider-unleash/%s (+https://github.com/evenh/terraform-provider-unleash) Terraform Plugin SDK/%s", version.ProviderVersion, meta.SDKVersionString())

	// TODO: Not hardcode
	auth := &api.UnsecureAuthentication{
		Email: "even.holthe@me.com",
	}

	client, err := api.NewClient(baseUrl, userAgent, auth)

	if err != nil {
		return nil, fmt.Errorf("could not configure Unleash REST client: %w", err)
	}

	return client, nil
}
