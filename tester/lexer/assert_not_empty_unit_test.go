package lexer

import (
	"reflect"
	"testing"
)

func Test_parsingAssertions_NotEmpty(t *testing.T) {
	data := `---assert_not_empty
		data.createCatalogModule.id
		data.createProduct.id
	`

	actual, err := parsingAssertions(data)
	if err != nil {
		t.Fail()
	}

	expected := []string{
		"data.createCatalogModule.id",
		"data.createProduct.id",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}
