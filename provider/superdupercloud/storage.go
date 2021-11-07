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

func storageResourceProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resCreateStorageContext,
		ReadContext:   resReadStorageContext,
		UpdateContext: resUpdateStorageContext,
		DeleteContext: resDeleteStorageContext,
		Importer:      importer(),

		Schema: map[string]*schema.Schema{
			"name":         nameSchema(),
			"display_name": displayNameSchema(),
			"size_mb": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntAtLeast(1),
				ExactlyOneOf:  []string{"network_id", "virtual_machine_id"},
				ConflictsWith: []string{"virtual_machine_id"},
			},
			"network_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateIP,
				RequiredWith:     []string{"network_id"},
			},
			"virtual_machine_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntAtLeast(1),
				RequiredWith:  []string{"mount_path"},
				ExactlyOneOf:  []string{"network_id", "virtual_machine_id"},
				ConflictsWith: []string{"network_id"},
			},
			"mount_path": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"virtual_machine_id"},
			},
		},
	}
}

func storageDataProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadStorageContext,

		Schema: map[string]*schema.Schema{
			"id":                 idSchema(),
			"name":               setComputed(nameSchema()),
			"display_name":       setComputed(displayNameSchema()),
			"size_mb":            {Type: schema.TypeInt, Computed: true},
			"network_id":         {Type: schema.TypeInt, Computed: true},
			"network_ip":         {Type: schema.TypeString, Computed: true},
			"virtual_machine_id": {Type: schema.TypeInt, Computed: true},
			"mount_path":         {Type: schema.TypeString, Computed: true},
		},
	}
}

func saveStorage(data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	storage := &entities.Storage{
		Name:             data.Get("name").(string),
		DisplayName:      toStrPtr(data.Get("display_name").(string)),
		SizeMB:           data.Get("size_mb").(int),
		NetworkID:        toIntPtr(data.Get("network_id").(int)),
		NetworkIP:        toStrPtr(data.Get("network_ip").(string)),
		VirtualMachineID: toIntPtr(data.Get("virtual_machine_id").(int)),
		MountPath:        toStrPtr(data.Get("mount_path").(string)),
	}

	if isInvalid(&d, storage) {
		return d
	}

	reqBody, err := json.Marshal(storage)
	if hasFailed(&d, err, "failed to marshall struct into JSON") {
		return d
	}

	var respBody []byte
	if data.Id() == "" {
		respBody, err = client.Create("/storage", reqBody)
	} else {
		respBody, err = client.Update("/storage/"+data.Id(), reqBody)
	}

	if hasFailed(&d, err, "request failed") {
		return d
	}

	return flattenStorage(d, data, respBody)
}

func resCreateStorageContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return saveStorage(data, meta)
}

func resReadStorageContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readStorage("/storage/"+data.Id(), data, meta)
}

func dataReadStorageContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readStorage("/storage/"+strconv.Itoa(data.Get("id").(int)), data, meta)
}

func readStorage(uri string, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	resp, err := client.Read(uri)
	if hasFailed(&d, err, "request failed") {
		return d
	}

	return flattenStorage(d, data, resp)
}

func resUpdateStorageContext(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return saveStorage(data, meta)
}

func resDeleteStorageContext(_ context.Context, data *schema.ResourceData, meta interface{}) (d diag.Diagnostics) {
	client := fromMeta(&d, meta)
	if client == nil {
		return d
	}

	if hasFailed(&d, client.Delete("/storage/"+data.Id()), "request failed") {
		return d
	}

	return d
}

func flattenStorage(d diag.Diagnostics, data *schema.ResourceData, body []byte) diag.Diagnostics {
	storage := &entities.Storage{}
	if err := json.Unmarshal(body, &storage); hasFailed(&d, err, "failed to unmarshal response into struct") {
		return d
	}

	data.SetId(strconv.FormatInt(storage.ID, 10))
	set(&d, data, "name", storage.Name)
	set(&d, data, "display_name", storage.DisplayName)
	set(&d, data, "size_mb", storage.SizeMB)
	set(&d, data, "network_id", storage.NetworkID)
	set(&d, data, "network_ip", storage.NetworkIP)
	set(&d, data, "virtual_machine_id", storage.VirtualMachineID)
	set(&d, data, "mount_path", storage.MountPath)

	return d
}
