package hal

import "fmt"

//
// Project
//

type Project struct {
	ResourceObject
}

func NewProject() *Project {
	return &Project{
		ResourceObject{
			Type: "Project",
		},
	}
}

func (res *Project) Id() int {
	return res.GetInt("id")
}

func (res *Project) Name() string {
	return res.GetString("name")
}

func (res *Project) Description() string {
	return res.GetString("description")
}

func (res *Project) GetWorkPackages(c *HalClient) (*Collection, error) {
	linkRes, err := res.GetLinkResource(c, "workPackages")
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
	resourceTypes["Project"] = func() Resource {
		return NewProject()
	}
}
