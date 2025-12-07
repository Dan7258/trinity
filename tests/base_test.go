package tests

import "testing"

func baseTest(t *testing.T) {
	if 1 == 2 {
		t.Error("Base test failed")
	}
}
