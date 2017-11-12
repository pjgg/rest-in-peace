package jsonAssert

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (self *jsonAssertSuite) TestAssertJsonEquals() {
	expected := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	actual := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsError() {
	expected := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	actual := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547311}`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.EqualError(err, "Json not equal. Fields [\"Time\"]: 1.294706395881547e+18 != 1.2947063958815473e+18")
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsTwoFieldsError() {
	expected := []byte(`{"Name":"Alice","Body":"Hell","Time":1294706395881547000}`)
	actual := []byte(`{"Name":"Alice","Body":"Hello","Time":12}`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.EqualError(err, "Json not equal. Fields [\"Body\"]: \"Hell\" != \"Hello\",[\"Time\"]: 1.294706395881547e+18 != 12")
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsDifferentOrder() {
	expected := []byte(`{"Body":"Hello","Name":"Alice","Time":1294706395881547000}`)
	actual := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsComplexJSON() {
	expected := []byte(`{ "id": "0001", "type": "donut", "name": "Cake", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)
	actual := []byte(`{ "id": "0001", "type": "donut", "name": "Cake", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsComplexJSONFieldError() {
	expected := []byte(`{ "id": "0001", "type": "donut", "name": "Cakee", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)
	actual := []byte(`{ "id": "0001", "type": "donut", "name": "Cake", "ppu": 0.55, "batters": { "batter": [ { "id": "1001", "type": "Regular" }, { "id": "1002", "type": "Chocolate" }, { "id": "1003", "type": "Blueberry" }, { "id": "1004", "type": "Devil's Food" } ] }, "topping": [ { "id": "5001", "type": "None" }, { "id": "5002", "type": "Glazed" }, { "id": "5005", "type": "Sugar" }, { "id": "5007", "type": "Powdered Sugar" }, { "id": "5006", "type": "Chocolate with Sprinkles" }, { "id": "5003", "type": "Chocolate" }, { "id": "5004", "type": "Maple" } ] }`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual); err != nil {
		self.EqualError(err, "Json not equal. Fields [\"name\"]: \"Cakee\" != \"Cake\"")
	}

}

func (self *jsonAssertSuite) TestAssertJsonEqualsWithIgnore() {
	expected := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
	actual := []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547311}`)

	if err := self.jsonAssert.AssertJsonEquals(expected, actual, "/Time"); err != nil {
		self.Fail("Unexpected error", err.Error())
	}

}

func TestJsonAssertSuite(t *testing.T) {
	testSuit := new(jsonAssertSuite)
	testSuit.jsonAssert = Instance()

	suite.Run(t, testSuit)
}

type jsonAssertSuite struct {
	suite.Suite
	jsonAssert JsonAssert
}
