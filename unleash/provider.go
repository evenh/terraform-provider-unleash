package unleash

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/meta"

	"github.com/evenh/terraform-provider-unleash/unleash/api"
	"github.com/evenh/terraform-provider-unleash/version"
)

// Environment Variables
const (
	UNLEASH_API_ENDPOINT  = "UNLEASH_API_ENDPOINT"
	UNLEASH_AUTH_USERNAME = "UNLEASH_AUTH_USERNAME"
	UNLEASH_AUTH_EMAIL    = "UNLEASH_AUTH_EMAIL"
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
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "Authentication mechanism to use for communicating with the Unleash API",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						auth_unsecure: {
							Type:     schema.TypeList,
							Optional: true,
							Required: false,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									auth_unsecure_email: {
										Type:        schema.TypeString,
										DefaultFunc: schema.EnvDefaultFunc(UNLEASH_AUTH_EMAIL, nil),
										Optional:    true,
									},
									auth_unsecure_username: {
										Type:         schema.TypeString,
										DefaultFunc:  schema.EnvDefaultFunc(UNLEASH_AUTH_USERNAME, nil),
										Optional:     true,
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

	providedAuth, err := expandAuthMechanism(d)
	if err != nil {
		return nil, fmt.Errorf("could not configure authentication mechanism: %w", err)
	}

	client, err := api.NewClient(baseUrl, userAgent, providedAuth)
	if err != nil {
		return nil, fmt.Errorf("could not configure Unleash REST client: %w", err)
	}

	return client, nil
}

func expandAuthMechanism(d *schema.ResourceData) (api.AuthMechanism, error) {
	var providedAuth api.AuthMechanism
	authMechanisms := map[string]map[string]interface{}{}

	if auth, ok := d.GetOk(auth); ok {
		// This is some hacky code for sure, PRs welcome
		for _, authMechanism := range auth.(*schema.Set).List() {
			for mechanismName, mechanismArgs := range authMechanism.(map[string]interface{}) {
				authMechanisms[mechanismName] = mechanismArgs.([]interface{})[0].(map[string]interface{})
			}
		}
	}

	// Handle unsecure (should be "insecure") auth scheme
	{
		var email, username string

		if m, ok := authMechanisms[auth_unsecure]; ok {
			// Read from config
			email = m[auth_unsecure_email].(string)
			username = m[auth_unsecure_username].(string)
		} else {
			// Read from env
			if v, ok := os.LookupEnv(UNLEASH_AUTH_EMAIL); ok {
				email = v
			}
			if v, ok := os.LookupEnv(UNLEASH_AUTH_USERNAME); ok {
				username = v
			}
		}

		if len(email) > 0 || len(username) > 0 {
			log.Println("[DEBUG] Using unsecure authentication")
			providedAuth = api.UnsecureAuthentication{
				Email:    email,
				Username: username,
			}
		}
	}

	// If no auth mechanism is specified
	if providedAuth == nil {
		return nil, fmt.Errorf("provider is missing authentication configuration")
	}

	return providedAuth, nil
}
