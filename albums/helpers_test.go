package albums

import "testing"

func assertExpectedError(errExpected bool, err error, t *testing.T) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
