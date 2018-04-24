package tester

import (
	"github.com/elgris/jsondiff"
)

var resolvDiff func(items []jsondiff.DiffItem) bool

func assertRequestContains(request, content map[string]interface{}) (bool, []byte) {
	diff := jsondiff.Compare(request, content)

	resolvDiff = func(items []jsondiff.DiffItem) bool {
		for _, i := range items {
			switch i.Resolution {
			case jsondiff.TypeAdded:
				return false
			case jsondiff.TypeNotEquals:
				return false
			case jsondiff.TypeDiff:
				return resolvDiff(i.ValueB.([]jsondiff.DiffItem))
			}
		}
		return true
	}

	if !resolvDiff(diff.Items()) {
		return false, jsondiff.Format(diff)
	}

	return true, nil
}
