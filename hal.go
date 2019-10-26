package hal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	duration "github.com/SpirentOrion/iso8601duration.v2"
)

var (
	resourceTypes = make(map[string]ResourceFactory)
)

type ResourceFactory = func() Resource

func NewResource(typeName string) Resource {
	factory, ok := resourceTypes[typeName]
	if ok {
		return factory()
	}
	// No specialized type, use generic ResourceObject
	return &ResourceObject{
		Type: typeName,
	}
}

type Link struct {
	Href       string      `json:"href"`
	Title      string      `json:"title,omitempty"`
	Templated  bool        `json:"templated,omitempty"`
	Method     string      `json:"method,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
	Identifier string      `json:"identifier,omitempty"`
}

type Resource interface {
	ResourceType() string
	GetLink(string) *Link
	IsError() *Error

	// Don't export this interface method.
	decodeHAL(map[string]json.RawMessage) error
}

//
// Generic Resource Object
//

type ResourceObject struct {
	Type  string          `json:"_type"`
	Links map[string]Link `json:"_links"`

	// Don't export these fields
	embedded map[string]interface{}
	fields   map[string]interface{}
}

func (res *ResourceObject) getField(field string) (interface{}, bool) {
	if res.fields != nil {
		if val, ok := res.fields[field]; ok {
			return val, true
		}
	}
	return nil, false
}

func (res *ResourceObject) HasField(field string) bool {
	_, ok := res.getField(field)
	return ok
}

func (res *ResourceObject) GetField(field string) interface{} {
	val, _ := res.getField(field)
	return val
}

func (res *ResourceObject) GetString(field string) string {
	val, ok := res.getField(field)
	if ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func (res *ResourceObject) GetInt(field string) int {
	val, ok := res.getField(field)
	if ok {
		switch n := val.(type) {
		case int:
			return n
		case float64:
			return int(n)
		case float32:
			return int(n)
		}
	}
	return 0
}

func (res *ResourceObject) GetDateTime(field string) (time.Time, error) {
	val := res.GetString(field)
	return time.Parse(time.RFC3339, val)
}

func (res *ResourceObject) GetDuration(field string) (time.Duration, error) {
	val := res.GetString(field)
	return duration.Parse(val)
}

func (res *ResourceObject) GetCreatedAt() *time.Time {
	if dt, err := res.GetDateTime("createdAt"); err == nil {
		return &dt
	}
	return nil
}

func (res *ResourceObject) GetUpdatedAt() *time.Time {
	if dt, err := res.GetDateTime("updatedAt"); err == nil {
		return &dt
	}
	return nil
}

func (res *ResourceObject) GetEmbeddedResource(name string, c *HalClient) Resource {
	if res.embedded != nil {
		if val, ok := res.embedded[name]; ok {
			if res, ok := val.(Resource); ok {
				return res
			}
		}
	}
	// Try loading from link
	if c != nil {
		linkRes, _ := res.GetLinkResource(c, name)
		return linkRes
	}
	return nil
}

func (res *ResourceObject) GetEmbeddedResourceList(name string) []Resource {
	if res.embedded != nil {
		if val, ok := res.embedded[name]; ok {
			if arr, ok := val.([]Resource); ok {
				return arr
			}
		}
	}
	return nil
}

func (res *ResourceObject) GetLink(name string) *Link {
	if link, ok := res.Links[name]; ok {
		return &link
	}
	return nil
}

func (res *ResourceObject) GetLinkResource(c *HalClient, name string) (Resource, error) {
	link := res.GetLink(name)
	if link == nil {
		return nil, errors.New("No Link")
	}
	// Request new page
	linkRes, err := c.LinkGet(link)
	if err != nil {
		// Failed to load page
		return nil, err
	}
	return linkRes, nil
}

func (res *ResourceObject) Delete(c *HalClient) error {
	link := res.GetLink("delete")
	if link == nil {
		return errors.New("No 'delete' Link")
	}
	// Delete this resource
	_, err := c.Delete(link.Href)
	if err != nil {
		// Failed to delete resource
		return err
	}
	return nil
}

func (res *ResourceObject) Update(c *HalClient) (Resource, error) {
	link := res.GetLink("updateImmediately")
	if link == nil {
		return nil, errors.New("No 'updateImmediately' Link")
	}
	// Patch this resource
	return c.Patch(link.Href, res)
}

func (res *ResourceObject) ResourceType() string {
	return res.Type
}

func (res *ResourceObject) IsError() *Error {
	return nil
}

func (res *ResourceObject) UnmarshalHAL(data []byte) error {
	return ResourceUnmarshal(res, data)
}

func (res *ResourceObject) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["_type"] = res.Type
	if res.Links != nil {
		m["_links"] = res.Links
	}
	if res.embedded != nil {
		m["_embedded"] = res.embedded
	}
	if res.fields != nil {
		for k, v := range res.fields {
			m[k] = v
		}
	}

	return json.Marshal(m)
}

func (res *ResourceObject) MarshalHAL() ([]byte, error) {
	return res.MarshalJSON()
}

func getFirstNonSpace(b json.RawMessage) byte {
	for _, c := range b {
		// http://godoc.org/encoding/json?file=scanner.go#isSpace
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		return c
	}
	// All white space, so just return a space.
	return ' '
}

func decodeEmbedded(buf json.RawMessage) (interface{}, error) {
	switch getFirstNonSpace(buf) {
	case '[':
		// Deocode as array of HAL Resources
		var raws []json.RawMessage
		if err := json.Unmarshal(buf, &raws); err != nil {
			return nil, err
		}
		arr := make([]Resource, 0, len(raws))
		for _, raw := range raws {
			if subRes, err := Unmarshal(raw); err != nil {
				return nil, err
			} else {
				arr = append(arr, subRes)
			}
		}
		return arr, nil
	case '{':
		// Deocode as HAL Resource
		if res, err := Unmarshal(buf); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	default:
		// Try decoding as basic JSON value.
		var val interface{}
		if err := json.Unmarshal(buf, &val); err != nil {
			return nil, err
		}
		return val, nil
	}
}

func (res *ResourceObject) decodeHAL(mData map[string]json.RawMessage) error {
	for key, val := range mData {
		switch key {
		case "_type":
		case "_links":
			// Unmarshal map of arrays of RawMessages
			var rawLinks map[string]json.RawMessage
			if err := json.Unmarshal(val, &rawLinks); err != nil {
				return err
			}
			// Unmarshal each Link or array of Links
			res.Links = make(map[string]Link)
			for key, val := range rawLinks {
				switch getFirstNonSpace(val) {
				case '{':
					var link Link
					if err := json.Unmarshal(val, &link); err != nil {
						log.Printf(" -- Unmarshal error: %s", err)
						return err
					}
					res.Links[key] = link
				case '[':
					// TODO: handle array of links
				default:
					log.Printf("---- Unknown Link value: [%s]", string(val))
				}
			}
		case "_embedded":
			// Unmarshal map of arrays of RawMessages
			var rawEmbedded map[string]json.RawMessage
			if err := json.Unmarshal(val, &rawEmbedded); err != nil {
				return err
			}
			// Unmarshal each embedded resource
			res.embedded = make(map[string]interface{})
			for key, val := range rawEmbedded {
				if val, err := decodeEmbedded(val); err != nil {
					return err
				} else {
					res.embedded[key] = val
				}
			}
		default:
			var field interface{}
			if err := json.Unmarshal(val, &field); err != nil {
				return err
			}
			if res.fields == nil {
				res.fields = make(map[string]interface{})
			}
			res.fields[key] = field
		}
	}
	return nil
}

func decodeResource(mData map[string]json.RawMessage, res Resource) (Resource, error) {
	// Decode resource type
	if typeRaw, ok := mData["_type"]; ok {
		var typeName string
		if err := json.Unmarshal(typeRaw, &typeName); err != nil {
			return nil, err
		}
		// Create new resource or load into existing resource.
		if res == nil {
			res = NewResource(typeName)
		} else if res.ResourceType() != typeName {
			// Expected resource type doesn't match decoded resource.
			if typeName == "Error" {
				// Decode `Error` resource and return as a Golang `error`
				resErr := NewError()
				if err := resErr.decodeHAL(mData); err != nil {
					// Failed to decode `Error` resource
					return nil, err
				}
				// Decoded `Error` return as a normal golang `error`
				return nil, resErr
			} else {
				// Programmer error?  They expected a different resource type.
				return nil, fmt.Errorf("Resource type mismatch: Expected '%s', got '%s'",
					res.ResourceType(), typeName)
			}
		}
	} else {
		return nil, fmt.Errorf("Missing '_type' field, unknown resource type.")
	}

	if err := res.decodeHAL(mData); err != nil {
		return nil, err
	}
	return res, nil
}

//
// Unmarshal and detect resource type
//
func Unmarshal(data []byte) (Resource, error) {
	// decode json
	var mData map[string]json.RawMessage
	if err := json.Unmarshal(data, &mData); err != nil {
		return nil, err
	}

	return decodeResource(mData, nil)
}

//
// Decode Resource from `io.Reader`
//
func Decode(r io.Reader) (Resource, error) {
	dec := json.NewDecoder(r)
	// decode json
	var mData map[string]json.RawMessage
	if err := dec.Decode(&mData); err != nil {
		return nil, err
	}

	return decodeResource(mData, nil)
}

//
// Unmarshal an expected Resource type
//
func ResourceUnmarshal(res Resource, data []byte) error {
	// decode json
	var mData map[string]json.RawMessage
	if err := json.Unmarshal(data, &mData); err != nil {
		return err
	}

	if _, err := decodeResource(mData, res); err != nil {
		return err
	}
	return nil
}

//
// Decode an expected Resource type from `io.Reader`
//
func ResourceDecode(res Resource, r io.Reader) error {
	dec := json.NewDecoder(r)
	// decode json
	var mData map[string]json.RawMessage
	if err := dec.Decode(&mData); err != nil {
		return err
	}

	if _, err := decodeResource(mData, res); err != nil {
		return err
	}
	return nil
}
