package example

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Third party service: Github connector code
type githubConnector struct {
	userName  string
	basicAuth string
	endpoint  string
}

type github interface {
	getCurrentUserInfo() map[string]interface{}
}

var gitHubConnectorSingleton sync.Once
var githubConnectorInstance githubConnector

func gitHubInstance(user, password, githubEndpoint string) *githubConnector {
	gitHubConnectorSingleton.Do(func() {
		githubConnectorInstance.userName = user
		githubConnectorInstance.endpoint = githubEndpoint
		githubConnectorInstance.basicAuth = base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	})

	return &githubConnectorInstance
}

func (git *githubConnector) getCurrentUserInfo() (currentUser map[string]interface{}) {
	fmt.Println(git.endpoint + "/users/" + git.userName)
	req, _ := newHTTPRequest("GET", git.endpoint+"/users/"+git.userName, nil, nil)
	req.Header.Add("Authorization", "basic "+git.basicAuth)
	if resp, err := makeHTTPQuery(req); err == nil {
		currentUser = make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&currentUser)
	}
	return
}

func newHTTPRequest(method, url string, body []byte, cookie *http.Cookie) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if nil != err {
		req = nil
	}

	if nil != cookie {
		req.AddCookie(cookie)
	}

	return req, err
}

func makeHTTPQuery(req *http.Request) (*http.Response, error) {

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}
