package hal

import (
	"fmt"
)

//
// Error
//

type Error struct {
	ResourceObject
}

// Golang Error interface
func (res *Error) Error() string {
	str := fmt.Sprintf("%s: %s", res.ErrorIdentifier(), res.Message())
	// Check for list of errors
	errors := res.GetEmbeddedResourceList("errors")
	for _, embed := range errors {
		if err, ok := embed.(*Error); ok {
			str += "\n" + err.Error()
		}
	}
	// Check for details
	details := res.GetEmbeddedResource("details", nil)
	if details != nil {
		if res, ok := details.(*ResourceObject); ok {
			for k, v := range res.fields {
				str += fmt.Sprintf("-- %s = %s", k, v)
			}
		} else {
			str += fmt.Sprintf("-- %+v", details)
		}
	}
	return str
}

func (res *Error) ErrorIdentifier() string {
	return res.GetString("errorIdentifier")
}

func (res *Error) Message() string {
	return res.GetString("message")
}

func (res *Error) ResourceType() string {
	return "Error"
}

func (res *Error) GetLink(string) *Link {
	return nil
}

func (res *Error) IsError() *Error {
	return res
}

func NewError() *Error {
	return &Error{
		ResourceObject{
			Type: "Error",
		},
	}
}

// Register Resource Factories
func init() {
	resourceTypes["Error"] = func() Resource {
		return NewError()
	}
}
