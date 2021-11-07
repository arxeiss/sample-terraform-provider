package superdupercloud

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

func set(d *diag.Diagnostics, data *schema.ResourceData, key string, value interface{}) {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return
	}
	var err error
	switch v := value.(type) {
	case *string:
		err = data.Set(key, *v)
	case *int:
		err = data.Set(key, *v)
	default:
		err = data.Set(key, value)
	}

	hasFailed(d, err, "failed to assign value %v into key %s", value, key)
}

func setComputed(base *schema.Schema) *schema.Schema {
	base.Required = false
	base.Optional = false
	base.Computed = true
	base.ForceNew = false
	base.ValidateDiagFunc = nil
	return base
}

func idSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntAtLeast(1),
	}
}

func nameSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: validateName,
	}
}

func displayNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDisplayName,
	}
}

func validateName(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{Severity: diag.Error, Summary: "expected type to be string", AttributePath: path},
		}
	}
	if err := entities.ValidateName(v); err != nil {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid name",
			Detail:        err.Error(),
			AttributePath: path,
		})
	}
	return ret
}

func validateDisplayName(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{Severity: diag.Error, Summary: "expected type to be string", AttributePath: path},
		}
	}
	if err := entities.ValidateDisplayName(&v); err != nil {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid display name",
			Detail:        err.Error(),
			AttributePath: path,
		})
	}
	return ret
}

func validateIP(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{Severity: diag.Error, Summary: "expected type to be string", AttributePath: path},
		}
	}
	if !entities.ValidIP(v) {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid IP address",
			AttributePath: path,
		})
	}
	return ret
}

func validateIPRange(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{Severity: diag.Error, Summary: "expected type to be string", AttributePath: path},
		}
	}
	if !entities.ValidIPRange(v) {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid IP range",
			Detail:        "IP range is not valid, use 8.8.8.8/24 format",
			AttributePath: path,
		})
	}
	return ret
}

func hasFailed(d *diag.Diagnostics, err error, summary string, args ...interface{}) bool {
	if err != nil {
		*d = append(*d, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf(summary, args...),
			Detail:   err.Error(),
		})
		return true
	}
	return false
}

func isInvalid(d *diag.Diagnostics, req entities.Validatable) bool {
	if err := req.Validate(); err != nil {
		*d = append(*d, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Entity validation failed",
			Detail:   err.Error(),
		})
		return true
	}
	return false
}

func toStrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func toIntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
