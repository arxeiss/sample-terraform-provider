package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/arxeiss/sample-terraform-provider/provider/superdupercloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider { //nolint:gocritic
			return superdupercloud.Provider()
		},
	})
}
