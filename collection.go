package hal

import (
	"fmt"
)

//
// Collection
//

type Collection struct {
	ResourceObject
}

func newCollection(typeName string) *Collection {
	return &Collection{
		ResourceObject{
			Type: typeName,
		},
	}
}

func NewCollection() *Collection {
	return newCollection("Collection")
}

func (res *Collection) Total() int {
	return res.GetInt("total")
}

func (res *Collection) Count() int {
	return res.GetInt("count")
}

func (res *Collection) IsPaginated() bool {
	return res.HasField("pageSize")
}

func (res *Collection) PageSize() int {
	return res.GetInt("pageSize")
}

func (res *Collection) Offset() int {
	return res.GetInt("offset")
}

func (res *Collection) getPage(c *HalClient, name string) (*Collection, error) {
	linkRes, err := res.GetLinkResource(c, name)
	if err != nil {
		return nil, err
	}
	// Make sure it is a Collection
	if col, ok := linkRes.(*Collection); ok {
		return col, nil
	}
	return nil, fmt.Errorf("Unknown resource type: %s", linkRes.ResourceType())
}

func (res *Collection) NextPage(c *HalClient) (*Collection, error) {
	return res.getPage(c, "nextByOffset")
}

func (res *Collection) PrevPage(c *HalClient) (*Collection, error) {
	return res.getPage(c, "previousByOffset")
}

func (res *Collection) Items() []Resource {
	return res.GetEmbeddedResourceList("elements")
}

// Register Resource Factories
func init() {
	resourceTypes["Collection"] = func() Resource {
		return NewCollection()
	}
	resourceTypes["WorkPackageCollection"] = func() Resource {
		return newCollection("WorkPackageCollection")
	}
}
