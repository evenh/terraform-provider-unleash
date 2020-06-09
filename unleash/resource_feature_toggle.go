package unleash

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// TODO: Add variants
func resourceFeatureToggle() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureToggleCreate,
		Read:   resourceFeatureToggleRead,
		Update: resourceFeatureToggleUpdate,
		Delete: resourceFeatureToggleDelete,
		Exists: resourceFeatureToggleExists,

		Importer: &schema.ResourceImporter{
			State: resourceFeatureToggleImport,
		},

		Schema: map[string]*schema.Schema{
			NAME: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique name for this feature toggle",
			},
			DESCRIPTION: {
				Type:        schema.TypeString,
				Required:    false,
				Description: "What this feature toggle represents",
			},
			ENABLED: {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     false,
				Description: "Whether this feature toggle should be enabled or not",
			},
			STRATEGIES: {
				// TODO: key is strategy name
				Type:        schema.TypeSet,
				Required:    false,
				Description: "Which strategies shall be applied to this feature toggle",
			},
		},
	}
}

func resourceFeatureToggleCreate(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureToggleRead(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureToggleUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureToggleDelete(d *schema.ResourceData, metaRaw interface{}) error {
	return nil
}

func resourceFeatureToggleExists(d *schema.ResourceData, metaRaw interface{}) (bool, error) {
	return false, nil
}

func resourceFeatureToggleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
