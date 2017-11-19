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

type HTTP_Method int

const (
	GET HTTP_Method = 1 + iota
	POST
	DELETE
	PUT
	PATCH
	HEAD
)

var HTTP_METHOD = [...]string{
	"GET",
	"POST",
	"DELETE",
	"PUT",
	"PATCH",
	"HEAD",
}

func (method HTTP_Method) String() string {
	return HTTP_METHOD[method-1]
}

type MockServer struct {
	stubWhen map[HTTP_Method]*stubReturn
	Port     int
}

type stubReturn struct {
	path       *regexp.Regexp
	thenReturn []byte
	status     int
}

type MockServerBehavior interface {
	When(httpMethod HTTP_Method, pathExpr string) StubReturnBehavior
	CleanStub()
}

type StubReturnBehavior interface {
	ThenReturn(thenReturn []byte, status int)
}

func Instance(port ...int) MockServerBehavior {
	mockServerSingleton.Do(func() {
		mockServer.stubWhen = make(map[HTTP_Method]*stubReturn)
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

func (self *MockServer) router(w http.ResponseWriter, r *http.Request) {
	var match bool
	for method, stub := range self.stubWhen {
		if r.Method == method.String() {
			if stub.path.MatchString(r.URL.Path) {
				w.WriteHeader(stub.status)
				fmt.Fprintf(w, string(stub.thenReturn))
				match = true
				break
			}
		}
	}

	if !match {
		panic(errors.New("No stub found for path" + r.URL.Path))
	}

}

func (self *MockServer) When(httpMethod HTTP_Method, pathExpr string) StubReturnBehavior {
	stubReturn := new(stubReturn)
	stubReturn.path = regexp.MustCompile(pathExpr)
	self.stubWhen[httpMethod] = stubReturn
	return stubReturn
}

func (self *MockServer) CleanStub() {
	self.stubWhen = make(map[HTTP_Method]*stubReturn)
}

func (self *stubReturn) ThenReturn(thenReturn []byte, status int) {
	self.thenReturn = thenReturn
	self.status = status
}
