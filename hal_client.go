package hal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

func (c *HalClient) newGet(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.base+path, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != nil {
		req.SetBasicAuth("apikey", *c.apiKey)
	}
	return req, nil
}

func (c *HalClient) newGetJSON(path string) (*http.Request, error) {
	req, err := c.newGet(path)
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

	return Decode(resp.Body)
}

func (c *HalClient) GetFile(path string) (io.Reader, error) {
	req, err := c.newGet(path)
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
	req, err := c.newGetJSON(path)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
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
