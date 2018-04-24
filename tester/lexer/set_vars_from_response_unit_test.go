package lexer

import (
	"reflect"
	"testing"

	"github.com/imega/graphql-tester/tester"
)

func Test_responseVars(t *testing.T) {
	vars := []byte(`


		@catalogid =      data.createCatalogModule.id


			@productid        = data.createProduct.id


	`)

	expected := tester.Vars{
		"@catalogid": ".data.createCatalogModule.id",
		"@productid": ".data.createProduct.id",
	}

	actual := responseVars(vars)

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}
