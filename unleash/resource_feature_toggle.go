package unleash

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/evenh/terraform-provider-unleash/unleash/api"
)

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
				Optional:    true,
				Description: "What this feature toggle represents",
			},
			ENABLED: {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Whether this feature toggle should be enabled or not",
			},
			/*STRATEGIES: {
				// TODO: key is strategy name
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Which strategies shall be applied to this feature toggle",
			},*/
		},
	}
}

func resourceFeatureToggleCreate(d *schema.ResourceData, metaRaw interface{}) error {
	client := metaRaw.(*api.Client)

	// TODO: Support more fields
	f := api.Feature{
		Name:        d.Get(NAME).(string),
		Description: d.Get(DESCRIPTION).(string),
		Enabled:     d.Get(ENABLED).(bool),
		Strategies:  []api.Strategy{{Name: "default", Parameters: api.ParameterMap{}}}, // TODO: Intentionally left to hardcoded default for now
	}

	err := client.CreateFeatureFlag(f)

	if err != nil {
		return fmt.Errorf("could not create feature toggle '%s': %w", d.Get(NAME), err)
	}

	return resourceFeatureToggleRead(d, metaRaw)
}

func resourceFeatureToggleRead(d *schema.ResourceData, metaRaw interface{}) error {
	client := metaRaw.(*api.Client)
	toggleName := d.Get(NAME).(string)

	f, err := client.FeatureFlagByName(toggleName)
	if err != nil {
		return fmt.Errorf("could not read feature toggle '%s': %w", d.Get(NAME), err)
	}

	d.SetId(f.Name)
	d.Set(NAME, f.Name)
	d.Set(DESCRIPTION, f.Description)
	d.Set(ENABLED, f.Enabled)

	return nil
}

func resourceFeatureToggleUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	client := metaRaw.(*api.Client)

	// TODO: Support more fields
	f := api.Feature{
		Name:        d.Get(NAME).(string),
		Description: d.Get(DESCRIPTION).(string),
		Enabled:     d.Get(ENABLED).(bool),
		Strategies:  []api.Strategy{{Name: "default", Parameters: api.ParameterMap{}}}, // TODO: Intentionally left to hardcoded default for now
	}

	err := client.UpdateFeatureFlag(d.Id(), f)
	if err != nil {
		return fmt.Errorf("could not update feature toggle '%s': %w", d.Get(NAME), err)
	}

	return resourceFeatureToggleRead(d, metaRaw)
}

func resourceFeatureToggleDelete(d *schema.ResourceData, metaRaw interface{}) error {
	client := metaRaw.(*api.Client)
	if err := client.DeleteFeatureFlag(d.Id()); err != nil {
		return fmt.Errorf("could not delete feature toggle '%s': %w", d.Get(NAME), err)
	}

	return nil
}

func resourceFeatureToggleExists(d *schema.ResourceData, metaRaw interface{}) (bool, error) {
	client := metaRaw.(*api.Client)
	toggleName := d.Get(NAME).(string)

	f, err := client.FeatureFlagByName(toggleName)

	return f != nil, err
}

func resourceFeatureToggleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	_ = d.Set(NAME, id)

	return []*schema.ResourceData{d}, nil
}
