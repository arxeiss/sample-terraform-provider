package superdupercloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	// Example Provider requires an API Token.
	// The Email is optional
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {Type: schema.TypeString, Required: true},
			"endpoint":  {Type: schema.TypeString, Optional: true},
		},

		ResourcesMap: map[string]*schema.Resource{
			"sdc_hdd": {
				Schema: map[string]*schema.Schema{
					"name":        {Type: schema.TypeString, Required: true},
					"description": {Type: schema.TypeString, Required: true},
				},
				CreateContext: func(c context.Context, rd *schema.ResourceData, i interface{}) diag.Diagnostics {
					log.Println("AAAAAAAAAAA ================================ ================================= AAAAAAAAAAAAAAAAAAAAAA")
					return nil
				},
				ReadContext: func(c context.Context, rd *schema.ResourceData, i interface{}) diag.Diagnostics {
					log.Println("AAAAAAAAAAA ================================ ================================= AAAAAAAAAAAAAAAAAAAAAA")
					return nil
				},
				UpdateContext: schema.NoopContext,
				DeleteContext: schema.NoopContext,
			},
		},
	}
}
