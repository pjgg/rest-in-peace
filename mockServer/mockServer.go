package mockServer

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var mockServerSingleton sync.Once
var mockServer MockServer

const portMax = 8989
const portMin = 8080

// HTTPMethod ... type GET, POST, DELETE, PUT, PATCH, HEAD
type HTTPMethod int

const (
	// GET ...
	GET HTTPMethod = 1 + iota
	// POST ...
	POST
	// DELETE ...
	DELETE
	// PUT ...
	PUT
	// PATCH ...
	PATCH
	// HEAD ...
	HEAD
)

var httpMethod = [...]string{
	"GET",
	"POST",
	"DELETE",
	"PUT",
	"PATCH",
	"HEAD",
}

func (method HTTPMethod) String() string {
	return httpMethod[method-1]
}

// MockServer represent the server that will contains all your stubs. When you are developing a microservice architecture you should talk with a lot of third party services, or other decouple service, so in order to test your solution you must mock all of these services. MockServer is the server that will mock those third party services.
type MockServer struct {
	stubWhen map[HTTPMethod]*stubReturn
	Port     int
}

type stubReturn struct {
	path       *regexp.Regexp
	thenReturn []byte
	status     int
	headers    map[string]string
}

// StubAction is the interface, the behavior of your mockServer. Don't forget clean your stubs(CleanStub) at the begining of each test.
type StubAction interface {
	When(httpMethod HTTPMethod, pathExpr string) StubReturn

	CleanStub()
}

// StubReturn will define the behavior of your "When" action.
type StubReturn interface {
	WithHeader(key string, value string) StubReturn

	ThenReturn(thenReturn []byte, status int)
}

// Instance return a singleton MockServer instance, this is why is important to clean your stubs before each test.
func Instance(port ...int) StubAction {
	mockServerSingleton.Do(func() {
		mockServer.stubWhen = make(map[HTTPMethod]*stubReturn)
		if len(port) > 0 {
			mockServer.Port = port[0]
		} else {
			rand.Seed(time.Now().Unix())
			mockServer.Port = rand.Intn(portMax-portMin) + portMin
		}
		http.HandleFunc("/", mockServer.router)
		go http.ListenAndServe(":"+strconv.Itoa(mockServer.Port), nil) // set listen port
		fmt.Println("Mock server up and running, listening over http://localhost:" + strconv.Itoa(mockServer.Port))
	})
	return &mockServer
}

func (mockServer *MockServer) router(w http.ResponseWriter, r *http.Request) {
	fullPath := mockServer.fullPath(r)
	var match bool
	for method, stub := range mockServer.stubWhen {
		if r.Method == method.String() {
			if stub.path.MatchString(fullPath) && mockServer.checkHeaders(r, stub.headers) {
				w = mockServer.buildResponse(w, stub)
				match = true
				break
			}
		}
	}

	if !match {
		panic(errors.New("No stub found for path" + fullPath))
	}

}

func (mockServer *MockServer) checkHeaders(r *http.Request, headers map[string]string) (match bool) {
	match = true
	if len(headers) > 0 {
		for expectedKey, expectedHeader := range headers {
			match = r.Header.Get(expectedKey) == expectedHeader
		}
	}

	return
}

func (mockServer *MockServer) buildResponse(w http.ResponseWriter, stub *stubReturn) http.ResponseWriter {
	w.WriteHeader(stub.status)
	fmt.Fprintf(w, string(stub.thenReturn))

	return w
}

func (mockServer *MockServer) fullPath(r *http.Request) (fullPath string) {
	if r.URL.RawQuery == "" {
		fullPath = r.URL.Path
	} else {
		fullPath = r.URL.Path + "?" + r.URL.RawQuery
	}

	return
}

// When ... define the precondition that should be achived in order to trigger the stubReturn. pathExpr must be a URL regular expression.
func (mockServer *MockServer) When(httpMethod HTTPMethod, pathExpr string) StubReturn {
	stubReturn := new(stubReturn)
	stubReturn.path = regexp.MustCompile(pathExpr)
	stubReturn.headers = make(map[string]string)
	mockServer.stubWhen[httpMethod] = stubReturn
	return stubReturn
}

// CleanStub ... cleans all stub defined previously
func (mockServer *MockServer) CleanStub() {
	mockServer.stubWhen = make(map[HTTPMethod]*stubReturn)
}

// ThenReturn, is the action that will be returned for a given precondition.
func (mockServer *stubReturn) ThenReturn(thenReturn []byte, status int) {
	mockServer.thenReturn = thenReturn
	mockServer.status = status
}

// WithHeader is used in order to add a header precondition that should be achived in order to trigger the stubReturn.
func (mockServer *stubReturn) WithHeader(key string, value string) StubReturn {
	mockServer.headers[key] = value
	return mockServer
}
