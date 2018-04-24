package tester

import (
	"testing"
)

func Test_replaceVarsToValuesInBody_String(t *testing.T) {

	vars := map[string]string{
		"@collectionID":  "value1value",
		"@collectionID2": "value2value",
		"@productID2":    "value3value",
		"@productID":     "value4value",
	}

	body := []byte(`
		{
			"variables": {
					"collectionID": "@collectionID",
					"itemID": "@collectionID2",
					"a": "@productID",
					"b": "@productID2"
			},
			"operationName": "ItemConnectionCreate"
		}
	`)

	b := replaceVarsToValuesInBody(vars, string(body))

	expected := `
		{
			"variables": {
					"collectionID": "value1value",
					"itemID": "value2value",
					"a": "value4value",
					"b": "value3value"
			},
			"operationName": "ItemConnectionCreate"
		}
	`

	if expected != b {
		t.Fail()
	}
}

func Test_replaceVarsToValuesInBody_Int(t *testing.T) {

	vars := map[string]string{
		"@collectionID":  "1",
		"@collectionID2": "2",
	}

	body := []byte(`
		{
			"variables": {
					"collectionID": @collectionID,
					"itemID": @collectionID2
			}
		}
	`)

	b := replaceVarsToValuesInBody(vars, string(body))

	expected := `
		{
			"variables": {
					"collectionID": 1,
					"itemID": 2
			}
		}
	`

	if expected != b {
		t.Fail()
	}
}
