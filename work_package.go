package hal

import "fmt"

//
// WorkPackage
//

type WorkPackage struct {
	ResourceObject
}

func NewWorkPackage() *WorkPackage {
	return &WorkPackage{
		ResourceObject{
			Type: "WorkPackage",
		},
	}
}

func (res *WorkPackage) Id() int {
	return res.GetInt("id")
}

func (res *WorkPackage) Description() *Formattable {
	f, err := DecodeFormattable(res.GetField("description"))
	if err != nil {
		return nil
	}
	return f
}

func (res *WorkPackage) Subject() string {
	return res.GetString("subject")
}

func (res *WorkPackage) GetAttachments(c *HalClient) (*Collection, error) {
	// Check if attachments are embedded
	val := res.GetEmbeddedResource("attachments")
	if col, ok := val.(*Collection); ok {
		return col, nil
	}
	// Load attachments
	linkRes, err := res.GetLinkResource(c, "attachments")
	if err != nil {
		return nil, err
	}
	// Make sure it is a Collection
	if col, ok := linkRes.(*Collection); ok {
		return col, nil
	}
	return nil, fmt.Errorf("Unknown resource type: %s", linkRes.ResourceType())
}

// Register Factories
func init() {
	resourceTypes["WorkPackage"] = func() Resource {
		return NewWorkPackage()
	}
}
