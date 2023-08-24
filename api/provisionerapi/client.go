package provisionerapi

import (
	"fmt"
	"time"
	"bytes"
	"reflect"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Client struct {
	apiEndpoint string
	apiKey		string
	httpClient	*http.Client
}

func New(apiEndpoint string, apiKey string) (*Client, diag.Diagnostics) {
	var diags diag.Diagnostics

	return &Client {
		apiEndpoint:	apiEndpoint,
		apiKey:			apiKey,
		httpClient:		&http.Client {
			Timeout: 100 * time.Second,
		},
	}, diags
}

func (client *Client) NewApiRequest(method string, path string, body interface{}, response interface{}) diag.Diagnostics {
	var jsonPayload []byte = nil
	var diags diag.Diagnostics
	var err error
	var httpRequest *http.Request

	if jsonPayload, err = json.Marshal(body); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"JSON serialization",
			Detail:		fmt.Sprintf("Failed to serialize JSON request, JSON payload: %s", jsonPayload),
		})
		
		return diags
	}

	apiCompletePathUrl := fmt.Sprintf("%s%s", client.apiEndpoint, path)

	if body != nil {
		httpRequest, err = http.NewRequest(method, apiCompletePathUrl, bytes.NewReader(jsonPayload))
	} else {
		httpRequest, err = http.NewRequest(method, apiCompletePathUrl, bytes.NewReader([]byte{}))
	}	

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"HTTP API request",
			Detail:		fmt.Sprintf("Failed to send HTTP request. Error: %s", err),
		})
		
		return diags
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("apikey", client.apiKey)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"HTTP API response",
			Detail:		fmt.Sprintf("HTTP response contains an error. Error: %s", err),
		})
		
		return diags
	}
	
	// close after NewApiRequest() is done
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		body, _ := ioutil.ReadAll(httpResponse.Body)
		diags = append(diags, diag.Diagnostic{
			Severity:	diag.Error,
			Summary:	"HTTP API response",
			Detail:		fmt.Sprintf("HTTP request failed. Status code: %d, response body: %s", httpResponse.StatusCode, string(body)),
		})
		
		return diags
	}

	if response != nil {
		responseBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:	diag.Error,
				Summary:	"HTTP API response",
				Detail:		fmt.Sprintf("HTTP request failed. Response: %s", err),
			})
			
			return diags
		}

		err = json.Unmarshal(responseBytes, response)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:	diag.Error,
				Summary:	"HTTP API response",
				Detail:		fmt.Sprintf("Failed to unmarshal HTTP response body as JSON to target struct of type '%s'\n, actual HTTP response payload is: '%s'\n", reflect.TypeOf(response).String(), string(responseBytes)),
			})
			
			return diags
		}
	}

	return nil
}