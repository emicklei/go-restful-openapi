package restfulspec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful"
)

func dummy(i *restful.Request, o *restful.Response) {}

type Sample struct {
	ID    string `swagger:"required"`
	Root  Item   `json:"root" description:"root desc"`
	Items []Item
}

type Item struct {
	ItemName string `json:"name"`
}

func asJSON(v interface{}) string {
	data, _ := json.MarshalIndent(v, " ", " ")
	return string(data)
}

func compareJson(t *testing.T, actualJsonAsString string, expectedJsonAsString string) bool {
	success := false
	var actualMap map[string]interface{}
	json.Unmarshal([]byte(actualJsonAsString), &actualMap)
	var expectedMap map[string]interface{}
	err := json.Unmarshal([]byte(expectedJsonAsString), &expectedMap)
	if err != nil {
		var actualArray []interface{}
		json.Unmarshal([]byte(actualJsonAsString), &actualArray)
		var expectedArray []interface{}
		err := json.Unmarshal([]byte(expectedJsonAsString), &expectedArray)
		success = reflect.DeepEqual(actualArray, expectedArray)
		if err != nil {
			t.Fatalf("Unparsable expected JSON: %s, actual: %v, expected: %v", err, actualJsonAsString, expectedJsonAsString)
		}
	} else {
		success = reflect.DeepEqual(actualMap, expectedMap)
	}
	if !success {
		t.Log("---- expected -----")
		t.Log(withLineNumbers(expectedJsonAsString))
		t.Log("---- actual -----")
		t.Log(withLineNumbers(actualJsonAsString))
		t.Log("---- raw -----")
		t.Log(actualJsonAsString)
		t.Error("there are differences")
		return false
	}
	return true
}
func withLineNumbers(content string) string {
	var buffer bytes.Buffer
	lines := strings.Split(content, "\n")
	for i, each := range lines {
		buffer.WriteString(fmt.Sprintf("%d:%s\n", i, each))
	}
	return buffer.String()
}
