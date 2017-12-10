package example

import (
	"testing"

	"github.com/pjgg/rest-in-peace/mockServer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const mockServerPort = 8080

func (testSuit *codeServiceSuite) SetupTest() {
	// Remember clean the stubs everytime that you run a test
	testSuit.mockServer.CleanStub()
}

func (testSuit *codeServiceSuite) TestGetGithubUserInfo() {

	// Mock response
	outboundJSON := []byte(`{ "login": "pjgg", "id": 3541131, "avatar_url": "https://avatars2.githubusercontent.com/u/3541131?v=4", "gravatar_id": "", "url": "https://api.github.com/users/pjgg", "html_url": "https://github.com/pjgg", "followers_url": "https://api.github.com/users/pjgg/followers", "following_url": "https://api.github.com/users/pjgg/following{/other_user}", "gists_url": "https://api.github.com/users/pjgg/gists{/gist_id}", "starred_url": "https://api.github.com/users/pjgg/starred{/owner}{/repo}", "subscriptions_url": "https://api.github.com/users/pjgg/subscriptions", "organizations_url": "https://api.github.com/users/pjgg/orgs", "repos_url": "https://api.github.com/users/pjgg/repos", "events_url": "https://api.github.com/users/pjgg/events{/privacy}", "received_events_url": "https://api.github.com/users/pjgg/received_events", "type": "User", "site_admin": false, "name": null, "company": null, "blog": "", "location": null, "email": null, "hireable": null, "bio": null, "public_repos": 18, "public_gists": 0, "followers": 5, "following": 1, "created_at": "2013-02-12T11:33:50Z", "updated_at": "2017-12-08T17:44:05Z", "private_gists": 0, "total_private_repos": 0, "owned_private_repos": 0, "disk_usage": 47035, "collaborators": 0, "two_factor_authentication": false, "plan": { "name": "free", "space": 976562499, "collaborators": 0, "private_repos": 0 } }`)

	// stubbing
	testSuit.mockServer.When(mockServer.GET, "/users/pjgg*").ThenReturn(outboundJSON, 200)

	// invoke
	userInfo := testSuit.codeService.getUserInfo()

	// asserts
	assert.Equal(testSuit.T(), "https://avatars2.githubusercontent.com/u/3541131?v=4", userInfo["avatar_url"].(string))
	assert.Equal(testSuit.T(), "pjgg", userInfo["login"].(string))
}

func TestCodeServiceSuite(t *testing.T) {
	testSuit := new(codeServiceSuite)
	testSuit.codeService = &codeService{
		// we don't really care about real user and password because this endpoint will be mocked!
		githubConnector: gitHubInstance("pjgg", "testPassword", "http://localhost:8080"),
	}
	testSuit.mockServer = mockServer.Instance(mockServerPort)

	suite.Run(t, testSuit)
}

type codeServiceSuite struct {
	suite.Suite
	codeService *codeService
	mockServer  mockServer.StubAction
}
