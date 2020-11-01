package unleash

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/evenh/terraform-provider-unleash/unleash/api"
)

// TODO: Add variants
func resourceFeatureToggle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFeatureToggleCreate,
		ReadContext:   resourceFeatureToggleRead,
		UpdateContext: resourceFeatureToggleUpdate,
		DeleteContext: resourceFeatureToggleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFeatureToggleImport,
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

func resourceFeatureToggleCreate(ctx context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	client := metaRaw.(*api.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// TODO: Support more fields
	f := api.Feature{
		Name:        d.Get(NAME).(string),
		Description: d.Get(DESCRIPTION).(string),
		Enabled:     d.Get(ENABLED).(bool),
		Strategies:  []api.Strategy{{Name: "default", Parameters: api.ParameterMap{}}}, // TODO: Intentionally left to hardcoded default for now
	}

	err := client.CreateFeatureFlag(f)

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not create feature toggle '%s': %w", d.Get(NAME), err))
	}

	resourceFeatureToggleRead(ctx, d, metaRaw)

	return diags
}

func resourceFeatureToggleRead(ctx context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	client := metaRaw.(*api.Client)
	toggleName := d.Get(NAME).(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	f, err := client.FeatureFlagByName(toggleName)
	if err != nil {
		// Signals that the flag doesn't exist upstream and should be removed from the state
		d.SetId("")
		return diag.FromErr(fmt.Errorf("could not read feature toggle '%s': %w", d.Get(NAME), err))
	}

	d.SetId(f.Name)
	d.Set(NAME, f.Name)
	d.Set(DESCRIPTION, f.Description)
	d.Set(ENABLED, f.Enabled)

	return diags
}

func resourceFeatureToggleUpdate(ctx context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
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
		return diag.FromErr(fmt.Errorf("could not update feature toggle '%s': %w", d.Get(NAME), err))
	}

	return resourceFeatureToggleRead(ctx, d, metaRaw)
}

func resourceFeatureToggleDelete(ctx context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	client := metaRaw.(*api.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if err := client.DeleteFeatureFlag(d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("could not delete feature toggle '%s': %w", d.Get(NAME), err))
	}

	return diags
}

func resourceFeatureToggleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	_ = d.Set(NAME, id)

	return []*schema.ResourceData{d}, nil
}
