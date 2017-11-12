// jsonAssert package. The aim of this package is provide json asserts methods, useful when you test a Rest API.
package jsonAssert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tonnerre/golang-pretty"
)

// jsonAssert interface define the main public behavior
type JsonAssert interface {
	AssertJsonEquals(expectedJson, actualJson []byte, ignorePaths ...string) error
}

type assertJsonImpl struct{}

// jsonAssert Instance is a factory in order to build JsonAsserts structs
func Instance() JsonAssert {
	return &assertJsonImpl{}
}

// AssertJsonEquals return error nil in case the inbound JSONs are the same, otherwise will return an error with the fields that are not the same.
// param 1: expectedJson, means the JSON that should be right in the assertion.
// param 2: actualJson, means the JSON that you are testing.
// param 3: (Optional) is an array of ignore paths. In other words, you will writedown all the json paths that must be ignored in your test assertion, ex /Time
func (self *assertJsonImpl) AssertJsonEquals(expectedJson, actualJson []byte, ignorePaths ...string) (err error) {
	var equal bool
	var msg []string
	var expected interface{}
	var actual interface{}

	if okExpected := json.Unmarshal(expectedJson, &expected); okExpected != nil {
		err = okExpected
		equal = false
		msg[0] = err.Error()
	}

	if okActual := json.Unmarshal(actualJson, &actual); okActual != nil {
		err = okActual
		equal = false
		msg[0] = err.Error()
	}

	expectedAsMap := expected.(map[string]interface{})
	actualAsMap := actual.(map[string]interface{})

	equal = reflect.DeepEqual(expectedAsMap, actualAsMap)
	if !equal {
		msg = pretty.Diff(expectedAsMap, actualAsMap)
	}

	if err == nil && !equal {
		if msgWithIgnorePathsExcluded := self.removeIgnorePaths(msg, ignorePaths); len(msgWithIgnorePathsExcluded) > 0 {
			err = fmt.Errorf("Json not equal. Fields " + strings.Join(msgWithIgnorePathsExcluded, ","))
		}
	}

	return
}

func (self *assertJsonImpl) removeIgnorePaths(msg []string, ignorePaths []string) (msgWithIgnorePathsExcluded []string) {
	var ignored bool
	for _, field := range msg {
		for _, ignorePath := range ignorePaths {
			path := self.fromPrettyMsgFormatToPath(field)
			if strings.Contains(strings.TrimSpace(path), strings.TrimSpace(ignorePath)) {
				ignored = true
				break

			}
		}
		if !ignored {
			msgWithIgnorePathsExcluded = append(msgWithIgnorePathsExcluded, field)
		}
		ignored = false
	}
	return
}

func (self *assertJsonImpl) fromPrettyMsgFormatToPath(msg string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(msg, "[\"", "/", -1), "\"]", "", -1), "[\"", "/", -1), "\"]", "", -1)
}
