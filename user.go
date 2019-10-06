package hal

import "fmt"

//
// User
//

type User struct {
	ResourceObject
}

func NewUser() *User {
	return &User{
		ResourceObject{
			Type: "User",
		},
	}
}

func (res *User) Id() int {
	return res.GetInt("id")
}

func (res *User) Name() string {
	return res.GetString("name")
}

func (res *User) FirstName() string {
	return res.GetString("firstName")
}

func (res *User) LastName() string {
	return res.GetString("lastName")
}

func (res *User) Login() string {
	return res.GetString("login")
}

func (res *User) Email() string {
	return res.GetString("email")
}

//
// UserPreferences
//

type UserPreferences struct {
	ResourceObject
}

func NewUserPreferences() *UserPreferences {
	return &UserPreferences{
		ResourceObject{
			Type: "UserPreferences",
		},
	}
}

func (res *UserPreferences) GetUser(c *HalClient) (*User, error) {
	linkRes, err := res.GetLinkResource(c, "user")
	if err != nil {
		return nil, err
	}
	// Make sure it is a User
	if u, ok := linkRes.(*User); ok {
		return u, nil
	}
	return nil, fmt.Errorf("Unknown resource type: %s", linkRes.ResourceType())
}

// Register Factories
func init() {
	resourceTypes["User"] = func() Resource {
		return NewUser()
	}
	resourceTypes["UserPreferences"] = func() Resource {
		return NewUserPreferences()
	}
}
