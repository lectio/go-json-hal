package hal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	// This API key is from a test user on a private test OpenProject instance.
	testAPIKey = "6d18cd23b6bd5e4c1531482b2a5eda171ce75cbc10c12d9c3370e898ef01ac69"
)

type testServer struct {
	server *httptest.Server

	client *HalClient

	router *http.ServeMux

	users map[string]string
}

func (ts *testServer) Close() {
	if ts.server != nil {
		ts.server.Close()
	}
}

func halErrorHandler(w http.ResponseWriter, status int, ident string, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"_type":"Error","errorIdentifier":"%s","messages":"%s"}`,
		ident, msg)
}

func (ts *testServer) checkAuth(w http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()
	if ok {
		// Check user & password
		if storedPass, ok := ts.users[user]; ok && storedPass == pass {
			// Valid username and password
			return
		}
	}
	halErrorHandler(w, http.StatusUnauthorized,
		"urn:openproject-org:api:v3:errors:Unauthenticated",
		"You need to be authenticated to access this resource.")
}

func (ts *testServer) addStatic(path string, data string, authRequired bool) {
	ts.router.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if authRequired {
			ts.checkAuth(w, req)
		}
		fmt.Fprintln(w, data)
	})
}

func newTestServer() *testServer {
	ts := &testServer{
		router: http.NewServeMux(),
		users:  make(map[string]string),
	}
	// Create test server
	server := httptest.NewServer(ts.router)
	// Connect client to test server
	ts.client = NewHalClient(server.URL)

	ts.users["apikey"] = testAPIKey

	// Default handler for "Not Found"
	ts.router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "{}")
	})

	// Add some json responses.
	ts.addStatic("/api/v3/configuration", `{
	"_type":"Configuration","maximumAttachmentFileSize":5242880,
	"perPageOptions":[20,100],"_links":{"self":{"href":"/api/v3/configuration"}}
}`, false)
	ts.addStatic("/api/v3/my_preferences", `{
	"_type":"UserPreferences","hideMail":true,"timeZone":null,
	"warnOnLeavingUnsaved":true,"commentSortDescending":false,"autoHidePopups":true,
	"_links":{
		"self":{"href":"/api/v3/my_preferences"},
		"user":{"href":"/api/v3/users/4","title":"test1 tester"},
		"updateImmediately":{"href":"/api/v3/my_preferences","method":"patch"}
	}
}`, true)

	return ts
}

func TestHalClient_Get(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	res, err := ts.client.Get("/api/v3/configuration")
	if err != nil {
		t.Errorf("HalClient failed to Get Hal resource: %v.", err)
	}
	if res == nil {
		t.Errorf("Resource expected.")
	}
}

func TestHalClient_LinkGet(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	// Load first resourcee
	res, err := ts.client.Get("/api/v3/configuration")
	if err != nil {
		t.Errorf("HalClient failed to Get Hal resource: %v.", err)
	}
	if res == nil {
		t.Errorf("Resource expected.")
	}

	// Get a Link from the resource
	link := res.GetLink("self")
	if link == nil {
		t.Errorf("Resource missing 'self' link.")
	}

	// Load linked resource.
	res2, err := ts.client.LinkGet(link)
	if err != nil {
		t.Errorf("HalClient failed to Get linked resource: %v.", err)
	}
	if res2 == nil {
		t.Errorf("Resource expected.")
	}
}

func TestHalClient_ApiKey(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	ts.client.SetAPIKey(testAPIKey)

	res, err := ts.client.Get("/api/v3/my_preferences")
	if err != nil {
		t.Errorf("HalClient failed to Get Hal resource: %v.", err)
		return
	}
	if res == nil {
		t.Errorf("Resource expected.")
	}
}

func TestHalClient_Unauthorized(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	res, err := ts.client.Get("/api/v3/my_preferences")
	if err != nil {
		if resErr, ok := err.(*Error); ok {
			if resErr.ErrorIdentifier() != "urn:openproject-org:api:v3:errors:Unauthenticated" {
				t.Errorf("Expected unauthorized response: %v.", err)
			}
		} else {
			t.Errorf("HalClient failed to Get Hal resource: %v.", err)
		}
	}
	if res != nil {
		t.Errorf("Expected unauthorized response: %v.", res)
	}
}
