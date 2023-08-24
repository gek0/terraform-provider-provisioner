package provisioner

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dapi "terraform-provider-provisioner/api/provisionerapi"
)

func resourceInstanceMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceMetadataCreate,
		ReadContext:   resourceInstanceMetadataRead,
		// update functionality is superfluous since every attribute change forces a new resource
		DeleteContext: resourceInstanceMetadataDelete,
		Schema: map[string]*schema.Schema{
			"ad_domain": &schema.Schema{
				Type:        		schema.TypeString,
				Required:	 		true,
				ForceNew:    		true,
				Description: 		"AD domain of the instance. Used for DC location",
				ValidateDiagFunc:	validateDomain,
			},
			"instance_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Hostname/name of the instance, without the domain part",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "DC location computed from ad_domain part",
			},
		},
	}
}

var instance_locations = map[string]string{
	"acme.local": "ac",
	"acme-replica.local": "me"
}

type deleteInstanceMetadataResponse struct {
	Status		string	 // success, error or partial
	Description	string
	Details		[]string // contains details when status == partial
}

func validateDomain(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	_, ok := instance_locations[v.(string)]

	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"Invalid AD domain name",
			Detail:		fmt.Sprintf("AD domain name %q not found in known acme domains", v.(string)),
		})

		return diags
	}

	return diags
}

func resourceInstanceMetadataCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	
	ad_domain := d.Get("ad_domain").(string)
	instance_name := d.Get("instance_name").(string)
	location := instance_locations[ad_domain]

	d.Set("location", location)
	d.SetId(instance_name + "." + ad_domain)

	return diags
}

func resourceInstanceMetadataRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if (d.Get("id") == "") {	
		// can be seen when enabling TF_LOG=<severity> option	
		fmt.Printf("[INFO] Instance metadata not found in state yet, ignoring.")
		d.SetId("")
	}

	// upstream API does not store any instance metadata so not returning aynthing else
	return diags
}

func resourceInstanceMetadataDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := fmt.Sprintf("/api/v1/vms/%s/%s", d.Get("location"), d.Get("instance_name"))
	deleteResponse := new(deleteInstanceMetadataResponse)

	// initialize API client and run DELETE request
	apiClient := m.(*dapi.Client)
	err := apiClient.NewApiRequest(http.MethodDelete, path, nil, deleteResponse)

	if err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness
	// can be seen when enabling TF_LOG=<severity> option
	log.Printf("[INFO] Instance metadata cleanup done. Status: %q. Details: %s", deleteResponse.Status, deleteResponse.Details)
	d.SetId("")

	return diags
}