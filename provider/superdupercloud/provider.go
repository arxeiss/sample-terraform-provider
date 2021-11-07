package superdupercloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type MetaContext struct {
	apiToken string
	endpoint string
	client   *http.Client
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {Type: schema.TypeString, Required: true},
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
		ConfigureContextFunc: func(c context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return &MetaContext{
				apiToken: rd.Get("api_token").(string),
				endpoint: rd.Get("endpoint").(string),
				client:   &http.Client{Timeout: 10 * time.Second},
			}, nil
		},
		ResourcesMap: map[string]*schema.Resource{
			"sdc_vm":      vmResourceProvider(),
			"sdc_storage": storageResourceProvider(),
			"sdc_network": networkResourceProvider(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"sdc_vm":      vmDataProvider(),
			"sdc_storage": storageDataProvider(),
			"sdc_network": networkDataProvider(),
		},
	}
}

func importer() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: basicStateImporter,
	}
}

func basicStateImporter(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	id, err := strconv.ParseInt(data.Id(), 10, 64)
	switch {
	case err != nil:
		return nil, errors.New("provided ID is not valid number")
	case id < 1:
		return nil, errors.New("ID must be number greater than 0")
	}

	return []*schema.ResourceData{data}, nil
}

func fromMeta(d *diag.Diagnostics, meta interface{}) *MetaContext {
	client, ok := meta.(*MetaContext)
	if !ok || client == nil {
		*d = append(*d, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable retrieve SuperDuperCloud client",
			Detail:   "Unable retrieve SuperDuperCloud client from meta",
		})
	}
	return client
}

func (mc *MetaContext) doRequest(method, uri string, body []byte, expectedCode int) ([]byte, error) {
	url := strings.TrimRight(mc.endpoint, "/") + "/" + strings.Trim(uri, "/")
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+mc.apiToken)
	req.Header.Set("Content-Type", "application/json")

	if len(body) > 0 {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	resp, err := mc.client.Do(req)
	if err != nil {
		return nil, err
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == expectedCode {
		return payload, nil
	}

	return nil, fmt.Errorf("request failed with code %d: %s", resp.StatusCode, string(payload))
}

func (mc *MetaContext) Create(uri string, body []byte) ([]byte, error) {
	return mc.doRequest("PUT", uri, body, http.StatusCreated)
}

func (mc *MetaContext) Read(uri string) ([]byte, error) {
	return mc.doRequest("GET", uri, nil, http.StatusOK)
}

func (mc *MetaContext) Update(uri string, body []byte) ([]byte, error) {
	return mc.doRequest("POST", uri, body, http.StatusOK)
}

func (mc *MetaContext) Delete(uri string) error {
	_, err := mc.doRequest("DELETE", uri, nil, http.StatusOK)
	return err
}
