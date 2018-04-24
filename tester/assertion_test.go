package tester

import (
	"testing"
)

func Test_assertRequestContains_Added_ReturnsTrue(t *testing.T) {

	request := map[string]interface{}{
		"id":    "1",
		"title": "title",
	}

	content := map[string]interface{}{
		"title": "title",
	}

	if ok, _ := assertRequestContains(request, content); !ok {
		t.Fail()
	}
}

func Test_assertRequestContains_Equal_ReturnsTrue(t *testing.T) {

	request := map[string]interface{}{
		"title": "title",
	}

	content := map[string]interface{}{
		"title": "title",
	}

	if ok, _ := assertRequestContains(request, content); !ok {
		t.Fail()
	}
}

func Test_assertRequestContains_Removed_ReturnsFalse(t *testing.T) {

	request := map[string]interface{}{
		"title": "title",
	}

	content := map[string]interface{}{
		"title":  "title",
		"second": "title",
	}

	if ok, _ := assertRequestContains(request, content); ok {
		t.Fail()
	}
}

func Test_assertRequestContains_NotEqual_ReturnsFalse(t *testing.T) {

	request := map[string]interface{}{
		"data": map[string]interface{}{
			"createCollection": map[string]interface{}{
				"id":          "4bcb0c0a25f76aa0bdc33381a4c262ef",
				"title":       "NewCollection",
				"description": "",
			},
		},
	}

	content := map[string]interface{}{
		"data": map[string]interface{}{
			"createCollection": map[string]interface{}{
				"title": "NewCollection1",
			},
		},
	}

	if ok, _ := assertRequestContains(request, content); ok {
		t.Fail()
	}
}
