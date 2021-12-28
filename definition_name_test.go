package restfulspec

import (
	"testing"
)

func TestGoLowerCamelCasedNameHandler_DefaultDefinitionName(t *testing.T) {
	testCases := []struct{
		name     string
		input    string
		excepted string
	}{
		{"go normal case","GoRestfulDefinition","goRestfulDefinition"},
		{"go ID case","IDForA","idForA"},
		{"go HTTP case","HTTPName","httpName"},
		{"go HTTPS case","HTTPSName","httpsName"},
		{"go HTTP with Started name case","HTTPStatus","httpStatus"},
	}

	for _, testCase := range testCases {
		output := GoLowerCamelCasedNameHandler(testCase.input)
		if output !=  testCase.excepted {
			t.Errorf("testing %s failed, expected %s, get %s",testCase.name, testCase.excepted, output)
		}
	}
}

func TestLowerCamelCasedNameHandler_DefaultDefinitionName(t *testing.T) {
	testCases := []struct{
		name     string
		input    string
		excepted string
	}{
		{"go normal case","GoRestfulDefinition","goRestfulDefinition"},
		{"go ID case","IDForA","iDForA"},
		{"go HTTP case","HTTPName","hTTPName"},
		{"go HTTPS case","HTTPSName","hTTPSName"},
		{"go HTTP with Started name case","HTTPStatus","hTTPStatus"},
	}

	for _, testCase := range testCases {
		output := LowerCamelCasedNameHandler(testCase.input)
		if output !=  testCase.excepted {
			t.Errorf("testing %s failed, expected %s, get %s",testCase.name, testCase.excepted, output)
		}
	}
}

func TestDefaultNameHandler_DefaultDefinitionName(t *testing.T) {
	testCases := []struct{
		name     string
		input    string
		excepted string
	}{
		{"go normal case","GoRestfulDefinition","GoRestfulDefinition"},
		{"go ID case","IDForA","IDForA"},
		{"go HTTP case","HTTPName","HTTPName"},
		{"go HTTPS case","HTTPSName","HTTPSName"},
		{"go HTTP with Started name case","HTTPStatus","HTTPStatus"},
	}

	for _, testCase := range testCases {
		output := DefaultNameHandler(testCase.input)
		if output !=  testCase.excepted {
			t.Errorf("testing %s failed, expected %s, get %s",testCase.name, testCase.excepted, output)
		}
	}
}

func TestLowerSnakeCasedNameHandler_DefaultDefinitionName(t *testing.T) {
	testCases := []struct{
		name     string
		input    string
		excepted string
	}{
		{"go normal case","GoRestfulDefinition","go_restful_definition"},
		{"go ID case","IDForA","i_d_for_a"},
	}

	for _, testCase := range testCases {
		output := LowerSnakeCasedNameHandler(testCase.input)
		if output !=  testCase.excepted {
			t.Errorf("testing %s failed, expected %s, get %s",testCase.name, testCase.excepted, output)
		}
	}
}
