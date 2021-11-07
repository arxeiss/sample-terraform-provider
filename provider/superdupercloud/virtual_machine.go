package superdupercloud

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func vmResourceProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resCreateVMContext,
		ReadContext:   resReadVMContext,
		UpdateContext: resUpdateVMContext,
		DeleteContext: resDeleteVMContext,
		Importer:      importer(),

		Schema: map[string]*schema.Schema{
			"name":         nameSchema(),
			"display_name": displayNameSchema(),
			"ram_size_mb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"network_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateIP,
				RequiredWith:     []string{"network_id"},
			},
			"public_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateIP,
			},
		},
	}
}
func vmDataProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadVMContext,

		Schema: map[string]*schema.Schema{
			"id":           idSchema(),
			"name":         setComputed(nameSchema()),
			"display_name": setComputed(displayNameSchema()),
			"ram_size_mb":  {Type: schema.TypeInt, Computed: true},
			"network_id":   {Type: schema.TypeInt, Computed: true},
			"network_ip":   {Type: schema.TypeString, Computed: true},
			"public_ip":    {Type: schema.TypeString, Computed: true},
		},
	}
}

func resCreateVMContext(_ context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	vm := &entities.VirtualMachine{
		Name:        data.Get("name").(string),
		DisplayName: toStrPtr(data.Get("display_name").(string)),
		RAMSizeMB:   data.Get("ram_size_mb").(int),
		NetworkID:   toIntPtr(data.Get("network_id").(int)),
		NetworkIP:   toStrPtr(data.Get("network_ip").(string)),
		PublicIP:    toStrPtr(data.Get("public_ip").(string)),
	}

	if isInvalid(&d, vm) {
		return d
	}

	reqBody, err := json.Marshal(vm)
	if hasFailed(&d, err, "failed to marshall struct into JSON") {
		return d
	}

	respBody, err := client.Create("/vm", reqBody)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	if err := json.Unmarshal(respBody, &vm); hasFailed(&d, err, "failed to unmarshal response into struct") {
		return d
	}

	data.SetId(strconv.FormatInt(vm.ID, 10))

	return d
}

func resReadVMContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readVM("/vm/"+data.Id(), data, meta)
}

func dataReadVMContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readVM("/vm/"+strconv.Itoa(data.Get("id").(int)), data, meta)
}

func readVM(uri string, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	resp, err := client.Read(uri)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	vm := &entities.VirtualMachine{}
	if err := json.Unmarshal(resp, &vm); hasFailed(&d, err, "failed to unmarshal response into struct") {
		return d
	}

	data.SetId(strconv.FormatInt(vm.ID, 10))
	set(&d, data, "name", vm.Name)
	set(&d, data, "display_name", vm.DisplayName)
	set(&d, data, "ram_size_mb", vm.RAMSizeMB)
	set(&d, data, "network_id", vm.NetworkID)
	set(&d, data, "network_ip", vm.NetworkIP)
	set(&d, data, "public_ip", vm.PublicIP)

	return d
}

func resUpdateVMContext(_ context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	vm := &entities.VirtualMachine{
		Name:        data.Get("name").(string),
		DisplayName: toStrPtr(data.Get("display_name").(string)),
		RAMSizeMB:   data.Get("ram_size_mb").(int),
		NetworkID:   toIntPtr(data.Get("network_id").(int)),
		NetworkIP:   toStrPtr(data.Get("network_ip").(string)),
		PublicIP:    toStrPtr(data.Get("public_ip").(string)),
	}

	if isInvalid(&d, vm) {
		return d
	}

	reqBody, err := json.Marshal(vm)
	if hasFailed(&d, err, "failed to marshall struct into JSON") {
		return d
	}

	_, err = client.Update("/vm/"+data.Id(), reqBody)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	return d
}

func resDeleteVMContext(_ context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	if hasFailed(&d, client.Delete("/vm/"+data.Id()), "request failed") {
		return d
	}

	return d
}
