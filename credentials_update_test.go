package main

import (
	"bytes"
	"encoding/json"
	"github.com/taskcluster/taskcluster-client-go/tcclient"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoutesTest struct {
	Routes
	t *testing.T
}

func TestCredentialsUpdate(t *testing.T) {
	newCreds := CredentialsUpdate{
		ClientId:    "newClientId",
		AccessToken: "newAccessToken",
		Certificate: "newCertificate",
	}

	body, err := json.Marshal(&newCreds)

	if err != nil {
		t.Fatal(err)
	}

	routes := NewRoutesTest(t)

	response := routes.request("POST", body)
	if response.Code != 405 {
		t.Errorf("Should return 405, but returned %d", response.Code)
	}

	response = routes.request("PUT", make([]byte, 0))
	if response.Code != 400 {
		t.Errorf("Should return 400, but returned %d", response.Code)
	}

	response = routes.request("PUT", body)
	if response.Code != 200 {
		content, _ := ioutil.ReadAll(response.Body)
		t.Fatal("Request error %d: %s", response.Code, string(content))
	}

	if routes.Credentials.ClientId != newCreds.ClientId {
		t.Errorf(
			"ClientId should be \"%s\", but got \"%s\"",
			newCreds.ClientId,
			routes.Credentials.ClientId,
		)
	}

	if routes.Credentials.AccessToken != newCreds.AccessToken {
		t.Errorf(
			"AccessToken should be \"%s\", but got \"%s\"",
			newCreds.AccessToken,
			routes.Credentials.AccessToken,
		)
	}

	if routes.Credentials.Certificate != newCreds.Certificate {
		t.Errorf(
			"Certificate should be \"%s\", but got \"%s\"",
			newCreds.Certificate,
			routes.Credentials.Certificate,
		)
	}
}

func (self *RoutesTest) request(method string, content []byte) (res *httptest.ResponseRecorder) {
	req, err := http.NewRequest(
		method,
		"http://localhost:8080/credentials",
		bytes.NewBuffer(content),
	)

	if err != nil {
		self.t.Fatal(err)
	}

	req.ContentLength = int64(len(content))
	res = httptest.NewRecorder()
	self.ServeHTTP(res, req)
	return
}

func NewRoutesTest(t *testing.T) *RoutesTest {
	return &RoutesTest{
		Routes: Routes{
			ConnectionData: tcclient.ConnectionData{
				Authenticate: true,
				Credentials: &tcclient.Credentials{
					ClientId:    "clientId",
					AccessToken: "accessToken",
					Certificate: "certificate",
				},
			},
		},
		t: t,
	}
}
