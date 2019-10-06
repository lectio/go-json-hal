package hal

import (
	"log"
	"testing"
)

func TestCollection_Unmarshal(t *testing.T) {
	s := `{
		"_type":"Collection","total":1,"count":1,
		"_embedded":{
			"elements":[
				{"_type":"Project","id":3,"identifier":"lectio",
					"name":"Lectio","description":"","createdAt":"2019-08-30T07:52:22Z",
					"updatedAt":"2019-08-31T11:08:11Z",
					"_links":{
						"self":{"href":"/api/v3/projects/3","title":"Lectio"},
						"createWorkPackage":{"href":"/api/v3/projects/3/work_packages/form","method":"post"},
						"createWorkPackageImmediate":{"href":"/api/v3/projects/3/work_packages","method":"post"},
						"workPackages":{"href":"/api/v3/projects/3/work_packages"},
						"categories":{"href":"/api/v3/projects/3/categories"},
						"versions":{"href":"/api/v3/projects/3/versions"},
						"types":{"href":"/api/v3/projects/3/types"}
					}
				}
			]
		},
		"_links":{
			"self":{"href":"/api/v3/projects"}
		}
}`
	res, err := Unmarshal([]byte(s))
	if err != nil {
		t.Errorf("Failed to parse Hal Collection %v.", err)
	}
	col, ok := res.(*Collection)
	if !ok {
		t.Errorf("Failed to cast Resource to Collection.")
	}
	log.Printf("Collection = %+v", col)
	// Get Embedded projects
	projects := col.Items()
	if projects == nil {
		t.Errorf("Collection missing 'elements'.")
	}
	if len(projects) != 1 {
		t.Errorf("Wrong number of projects: %d != 1", len(projects))
	}
	log.Printf("Project = %+v", projects[0])
}

func TestError_Unmarshal(t *testing.T) {
	s := `{"_type": "Error",
	"errorIdentifier": "urn:openproject-org:api:v3:errors:InternalServerError",
	"message": "An internal server error occured. This is not your fault." }`
	res, err := Unmarshal([]byte(s))
	if err != nil {
		t.Errorf("Failed to parse Hal Error %v.", err)
	}
	resErr := res.IsError()
	if resErr == nil {
		t.Errorf("Hal resource isn't an Error object.")
	}
	log.Printf("Error = %+v", resErr)
}
