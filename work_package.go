package hal

import "time"

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

func (res *WorkPackage) GetAttachments(c *HalClient) *Collection {
	// Get embedded attachments or load from a link
	val := res.GetEmbeddedResource("attachments", c)
	if col, ok := val.(*Collection); ok {
		return col
	}
	return nil
}

func (res *WorkPackage) GetAuthor(c *HalClient) *User {
	// Get embedded author or load from a link
	val := res.GetEmbeddedResource("author", c)
	if u, ok := val.(*User); ok {
		return u
	}
	return nil
}

func (res *WorkPackage) GetResponsible(c *HalClient) *User {
	// Get embedded responsible or load from a link
	val := res.GetEmbeddedResource("responsible", c)
	if u, ok := val.(*User); ok {
		return u
	}
	return nil
}

func (res *WorkPackage) GetAssignee(c *HalClient) *User {
	// Get embedded assignee or load from a link
	val := res.GetEmbeddedResource("assignee", c)
	if u, ok := val.(*User); ok {
		return u
	}
	return nil
}

func (res *WorkPackage) AddTimeEntry(c *HalClient, te *TimeEntry) (Resource, error) {
	if l := res.GetLink("project"); l != nil {
		te.AddLink("project", *l)
	}
	if l := res.GetLink("self"); l != nil {
		te.AddLink("workPackage", *l)
	}
	// Make sure 'spentOn' is set.
	if !te.HasField("spentOn") {
		te.SetSpentOn(time.Now())
	}

	return c.Post("/api/v3/time_entries", te)
}

// Register Factories
func init() {
	resourceTypes["WorkPackage"] = func() Resource {
		return NewWorkPackage()
	}
}
