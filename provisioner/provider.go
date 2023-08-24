package provisioner

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dapi "terraform-provider-provisioner/api/provisionerapi"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": &schema.Schema{
				Type:			schema.TypeString,
				Optional:		true,
				Description:	"API URL that the provider uses to utilize its functionalities. Optional",
				DefaultFunc:	schema.EnvDefaultFunc("PROVISIONER_API_ENDPOINT", "https://vm-deploy.acme.com"),
			},
			"api_key": &schema.Schema{
				Type:			schema.TypeString,
				Required:		true,
				Description:	"API key of the user (administrator access or less for specific virtual machines)",
				DefaultFunc:	schema.EnvDefaultFunc("PROVISIONER_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"provisioner_instance_metadata": resourceInstanceMetadata(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
		
	api_endpoint := d.Get("api_endpoint").(string)
	api_key := d.Get("api_key").(string)
	provisionerClient, err := dapi.New(api_endpoint, api_key)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"Unable to create Provisioner API client",
			Detail:		"Unexpected error during client intialization",
		})

		return nil, diags		
	}

	return provisionerClient, diags
}