package main

import (
	"log"

	"github.com/evenh/terraform-provider-unleash/unleash/api"
)

func main() {
	/*plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return unleash.Provider()
		},
	})*/

	client, err := api.NewClient("http://localhost:4242/api", "terraform-provider-unleash", &api.UnsecureAuthentication{Email: "even.holthe@me.com"})

	if err != nil {
		log.Fatal(err)
	}

	features, err := client.ListFeatureFlags()

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range features {
		log.Printf("Feature[name=%s, description=%s, enabled=%s]", f.Name, f.Description, f.Enabled)
	}
}
