package superdupercloud

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func networkResourceProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resCreateNetworkContext,
		ReadContext:   resReadNetworkContext,
		UpdateContext: resUpdateNetworkContext,
		DeleteContext: resDeleteNetworkContext,
		Importer:      importer(),

		Schema: map[string]*schema.Schema{
			"name":         nameSchema(),
			"display_name": displayNameSchema(),
			"ip_range": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIPRange,
			},
			"use_dhcp": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func networkDataProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadNetworkContext,

		Schema: map[string]*schema.Schema{
			"id":           idSchema(),
			"name":         setComputed(nameSchema()),
			"display_name": setComputed(displayNameSchema()),
			"ip_range":     {Type: schema.TypeString, Computed: true},
			"use_dhcp":     {Type: schema.TypeBool, Computed: true},
		},
	}
}

func resCreateNetworkContext(ctx context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	network := &entities.Network{
		Name:        data.Get("name").(string),
		DisplayName: toStrPtr(data.Get("display_name").(string)),
		IPRange:     data.Get("ip_range").(string),
		UseDHCP:     data.Get("use_dhcp").(bool),
	}

	if isInvalid(&d, network) {
		return d
	}

	reqBody, err := json.Marshal(network)
	if hasFailed(&d, err, "failed to marshall struct into JSON") {
		return d
	}

	respBody, err := client.Create("/network", reqBody)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	if err := json.Unmarshal(respBody, &network); hasFailed(&d, err, "failed to unmarshal response into struct") {
		return d
	}

	data.SetId(strconv.FormatInt(network.ID, 10))

	return d
}

func resReadNetworkContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readNetwork("/network/"+data.Id(), data, meta)
}

func dataReadNetworkContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readNetwork("/network/"+strconv.Itoa(data.Get("id").(int)), data, meta)
}

func readNetwork(uri string, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	resp, err := client.Read(uri)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	network := &entities.Network{}
	if err := json.Unmarshal(resp, &network); hasFailed(&d, err, "failed to unmarshal response into struct") {
		return d
	}

	data.SetId(strconv.FormatInt(network.ID, 10))
	set(&d, data, "name", network.Name)
	set(&d, data, "display_name", network.DisplayName)
	set(&d, data, "ip_range", network.IPRange)
	set(&d, data, "use_dhcp", network.UseDHCP)

	return d
}

func resUpdateNetworkContext(ctx context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	network := &entities.Network{
		Name:        data.Get("name").(string),
		DisplayName: toStrPtr(data.Get("display_name").(string)),
		IPRange:     data.Get("ip_range").(string),
		UseDHCP:     data.Get("use_dhcp").(bool),
	}

	if isInvalid(&d, network) {
		return d
	}

	reqBody, err := json.Marshal(network)
	if hasFailed(&d, err, "failed to marshall struct into JSON") {
		return d
	}

	_, err = client.Update("/network/"+data.Id(), reqBody)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	return d
}

func resDeleteNetworkContext(ctx context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	if hasFailed(&d, client.Delete("/network/"+data.Id()), "request failed") {
		return d
	}

	return d
}
