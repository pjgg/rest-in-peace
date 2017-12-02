package mockServer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pjgg/rest-in-peace/jsonAssert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const mockServerPort = 8080

func (self *mockServerSuite) SetupTest() {
	// Remember clean the stubs everytime that you run a test
	self.mockServer.CleanStub()
}

func (self *mockServerSuite) TestHelloWorld() {
	var err error
	outboundJson := []byte(`{"key":"hello","value":"world"}`)

	self.mockServer.When(GET, "/v1/service/hello*").ThenReturn(outboundJson, 200)

	req, _ := newHTTPRequest("GET", "http://localhost:8080/v1/service/hello/", nil, nil)
	if resp, err := makeHTTPQuery(req); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(self.T(), self.jsonAssert.AssertJsonEquals(outboundJson, data), "Unexpected error")
		assert.Equal(self.T(), 200, resp.StatusCode)
	}

	if err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func (self *mockServerSuite) TestComplexResponse() {
	var err error
	outboundJson := []byte(`{ "id": "0001", "type": "donut", "name": "Cake", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)

	self.mockServer.When(GET, "/v1/service/hello\\?param=\\d+").ThenReturn(outboundJson, 200)

	req, _ := newHTTPRequest("GET", "http://localhost:8080/v1/service/hello?param=10", nil, nil)
	if resp, err := makeHTTPQuery(req); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(self.T(), self.jsonAssert.AssertJsonEquals(outboundJson, data), "Unexpected error")
		assert.Equal(self.T(), 200, resp.StatusCode)
	}

	if err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func (self *mockServerSuite) TestRequestWithMultipleQueryParams() {
	var err error
	outboundJson := []byte(`{ "id": "0001", "type": "donut", "name": "Cake", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)

	self.mockServer.When(GET, "/v1/service/hello\\?param=\\d+&param_two=\\d+").ThenReturn(outboundJson, 200)

	req, _ := newHTTPRequest("GET", "http://localhost:8080/v1/service/hello?param=10&param_two=11", nil, nil)
	if resp, err := makeHTTPQuery(req); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)

		assert.Nil(self.T(), self.jsonAssert.AssertJsonEquals(outboundJson, data), "Unexpected error")
		assert.Equal(self.T(), 200, resp.StatusCode)
	}

	if err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func TestMockServerSuite(t *testing.T) {
	testSuit := new(mockServerSuite)
	testSuit.mockServer = Instance(mockServerPort)
	testSuit.jsonAssert = jsonAssert.Instance()
	suite.Run(t, testSuit)
}

type mockServerSuite struct {
	suite.Suite
	mockServer MockServerBehavior
	jsonAssert jsonAssert.JsonAssert
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
