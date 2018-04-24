package tester

import (
	"reflect"
	"testing"
)

func TestPadding_WithoutWrap(t *testing.T) {
	actual := padding(
		"qwertyuiopasdfghjklzxcvbnm",
		"ok",
		false,
	)

	expected := "qwertyuiopasdfghjklzxcvbnm..........................................ok"

	if expected != actual {
		t.Fail()
	}
}

func TestPadding_WithoutWrap_ReturnsCutString(t *testing.T) {
	actual := padding(
		"qwertyuiopasdfghjklzxcvbnm%qwertyuiopasdfghjklzxcvbnm%qwertyuiopasdfghjklzxcvbnm",
		"ok",
		false,
	)

	expected := "qwertyuiopasdfghjklzxcvbnm%qwertyuiopasdfghjklzxcvbnm%qwertyuiopa...ok"

	if expected != actual {
		t.Fail()
	}
}

func TestPadding_WithWrap(t *testing.T) {
	actual := padding(
		"qwertyuiopasdfghjklzxcvbnm",
		"ok",
		true,
	)

	expected := "qwertyuiopasdfghjklzxcvbnm..........................................ok"

	if expected != actual {
		t.Fail()
	}
}

func TestPadding_WithWrap_ReturnsWrappedString(t *testing.T) {
	actual := padding(
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque suscipit erat non sagittis egestas. In non ultrices leo, ac gravida dui. Vestibulum mollis felis vitae orci venenatis, eget elementum neque feugiat",
		"ok",
		true,
	)

	expected := "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n" +
		"Pellentesque suscipit erat non sagittis egestas. In non ultrices\n" +
		"leo, ac gravida dui. Vestibulum mollis felis vitae orci\n" +
		"venenatis, eget elementum neque feugiat.............................ok"

	if expected != actual {
		t.Fail()
	}
}

func TestWordWrap(t *testing.T) {
	actual := wordWrap(
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque suscipit erat non sagittis egestas. In non ultrices leo, ac gravida dui. Vestibulum mollis felis vitae orci venenatis, eget elementum neque feugiat.",
		70,
	)

	expected := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque",
		"suscipit erat non sagittis egestas. In non ultrices leo, ac gravida",
		"dui. Vestibulum mollis felis vitae orci venenatis, eget elementum",
		"neque feugiat.",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}
