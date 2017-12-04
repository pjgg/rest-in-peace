# Rest-in-peace 

[![Build Status](https://travis-ci.org/pjgg/rest-in-peace.svg?branch=master)](https://travis-ci.org/pjgg/rest-in-peace)
[![codecov](https://codecov.io/gh/pjgg/rest-in-peace/branch/master/graph/badge.svg)](https://codecov.io/gh/pjgg/rest-in-peace)

The aim of this project is to mock the microservices responses that your service depends on. By this way you will not need that third party services will be up and running in order to test and develop your microservice, because each time that you request some data to this APIs you will get a mock response. 

![Basic Usage](https://github.com/pjgg/rest-in-peace/blob/master/basicDiagram.png "Basic usage")

## How to use

1. Define MockServer port

```golang
const mockServerPort = 8080 
```

2. Get a singleton instance of your MockServer

```golang
mockServer = Instance(mockServerPort)
```

3. Clean your stubs each time that you run a test

```golang
mockServer.CleanStub()
```

4. define your test stubs

```golang
outboundJson := []byte(`{"key":"hello","value":"world"}`)
mockServer.When(GET, "/v1/service/hello*").ThenReturn(outboundJson, 200)
```

## Testify Integration example

This library needs from other libraries in order to build and orchestrate your unit/integration test. I will provide an integration example using testify. 

1. Define your test suit struct

```golang
type mockServerSuite struct {
	suite.Suite
	mockServer MockServerBehavior
	jsonAssert jsonAssert.JsonAssert
}
```

2. Setup testify main function and instanciate your mockServer

```golang
const mockServerPort = 8080
func TestMockServerSuite(t *testing.T) {
	testSuit := new(mockServerSuite)
	testSuit.mockServer = Instance(mockServerPort)
	testSuit.jsonAssert = jsonAssert.Instance()
	suite.Run(t, testSuit)
}
```

3. Define your test setup, the code that must be run before each test.

```golang
func (self *mockServerSuite) SetupTest() {
	// Remember clean the stubs everytime that you run a test
	self.mockServer.CleanStub()
}
```

4. Write down your unit test

```golang
func (self *mockServerSuite) TestHelloWorld() {
	var err error
	outboundJson := []byte(`{"key":"hello","value":"world"}`)

	self.mockServer.When(GET, "/v1/service/hello*").ThenReturn(outboundJson, 200)

    // Your service as myService.method() or a simple httpRequest
    req, _ := newHTTPRequest("GET", "http://localhost:8080/v1/service/hello/", nil, nil)
    
    // assert your response.
	if resp, err := makeHTTPQuery(req); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)

        // You could use restInPeace jsonAssert or testify asserts or your favorite assert lib
		assert.Nil(self.T(), self.jsonAssert.AssertJsonEquals(outboundJson, data), "Unexpected error")
		assert.Equal(self.T(), 200, resp.StatusCode)
	}

	if err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}
```