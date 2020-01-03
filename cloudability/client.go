package cloudability

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"path"
)

const (
	// endpoints
	api_v1_url = "https://app.cloudability.com"
	api_v3_url = "https://api.cloudability.com"
)

// Client is a Cloudability http client
type cloudabilityClient struct {
	BusinessMappings *businessMappingsEndpoint
	// Users *UsersEndpoint
	// Vendors *VendorsEndpoint
	// Views *ViewsEndpoint
}

func NewCloudabilityClient(apikey string) *cloudabilityClient {
	c := &cloudabilityClient{}
	c.BusinessMappings = newBusinessMappingsEndpoint(apikey)
	return c
}

type cloudabilityV3Endpoint struct {
	*http.Client
	BaseURL *url.URL
	EndpointPath string
	UserAgent string
	apikey string
}

type cloudabilityV1Endpoint struct {
	*cloudabilityV3Endpoint
}

func newCloudabilityV3Endpoint(apikey string) *cloudabilityV3Endpoint {
	e := &cloudabilityV3Endpoint{
		Client: &http.Client{Timeout: 10 * time.Second},
		UserAgent: "cloudability-sdk-go",
		apikey: apikey,
	}
	e.BaseURL, _ = url.Parse(api_v3_url)
	return e
}

func newCloudabilityV1Endpoint(apikey string) *cloudabilityV1Endpoint {
	e := &cloudabilityV1Endpoint{newCloudabilityV3Endpoint(apikey)}
	e.BaseURL, _ = url.Parse(api_v1_url)
	return e
}

func (e* cloudabilityV3Endpoint) get(endpoint string, result interface{}) error {
	endpointPath := path.Join(e.EndpointPath, endpoint)
	req, err := e.newVRequest("GET", endpointPath, nil)
	if err != nil {
		return err
	}
	_, err = e.execRequest(req, &result)
	return err
}

func (ce* cloudabilityV3Endpoint) execRequest(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := ce.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(bodyBytes))
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		log.Fatal(err)
	}
	return resp, nil
}

func (ce* cloudabilityV3Endpoint) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := ce.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", ce.UserAgent)
	return req, nil
}

func (ce cloudabilityV1Endpoint) newVRequest(method string, path string, body interface{}) (*http.Request, error) {
	req, err := ce.newRequest(method, path,body)
	q := req.URL.Query()
	q.Add("auth_token", ce.apikey)
	req.URL.RawQuery = q.Encode()
	return req, err
}

func (ce cloudabilityV3Endpoint) newVRequest(method string, path string, body interface{}) (*http.Request, error) {
	req, err := ce.newRequest(method, path, body)
	req.SetBasicAuth(ce.apikey, "")
	return req, err
}