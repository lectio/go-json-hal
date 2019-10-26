package hal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HalClient struct {
	http.Client

	base string

	// API Key Auth
	apiKey *string
}

func NewHalClient(base string) *HalClient {
	return &HalClient{
		base: base,
	}
}

func (c *HalClient) SetAPIKey(key string) {
	c.apiKey = &key
}

func (c *HalClient) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return nil, err
	}
	if c.apiKey != nil {
		req.SetBasicAuth("apikey", *c.apiKey)
	}
	return req, nil
}

func (c *HalClient) newRequestJSON(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *HalClient) doRequest(req *http.Request) (Resource, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := Decode(resp.Body)
	if err != nil {
		// HTTP Error
		return nil, err
	}
	// Convert HAL errors to golang error
	if err := res.IsError(); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *HalClient) GetFile(path string) (io.Reader, error) {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *HalClient) Get(path string) (Resource, error) {
	req, err := c.newRequestJSON("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

func (c *HalClient) Delete(path string) (*http.Response, error) {
	req, err := c.newRequestJSON("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HalClient) Post(path string, res Resource) (Resource, error) {
	// encode resource as JSON for Post body.
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	if req, err := c.newRequestJSON("POST", path, bytes.NewBuffer(body)); err != nil {
		return nil, err
	} else {
		return c.doRequest(req)
	}
}

func (c *HalClient) Patch(path string, res Resource) (Resource, error) {
	// encode resource as JSON for Post body.
	body, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	if req, err := c.newRequestJSON("PATCH", path, bytes.NewBuffer(body)); err != nil {
		return nil, err
	} else {
		return c.doRequest(req)
	}
}

func (c *HalClient) GetCollection(path string) (*Collection, error) {
	res, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	col, ok := res.(*Collection)
	if !ok {
		return nil, fmt.Errorf("Invalid resource type: %s", res.ResourceType())
	}

	return col, nil
}

func (c *HalClient) GetFilteredCollection(path string, filters *Filters) (*Collection, error) {
	if f := filters.String(); f != "" {
		path += "?filters=" + url.QueryEscape(f)
	}
	return c.GetCollection(path)
}

func (c *HalClient) LinkGet(link *Link) (Resource, error) {
	if link == nil {
		return nil, errors.New("nil Link")
	}
	return c.Get(link.Href)
}

func (c *HalClient) LinkGetFile(link *Link) (io.Reader, error) {
	if link == nil {
		return nil, errors.New("nil Link")
	}
	return c.GetFile(link.Href)
}
