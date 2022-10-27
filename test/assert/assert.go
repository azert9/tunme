package assert

import "testing"

func NotEqual[T comparable](t *testing.T, a T, b T) {

	if a == b {
		t.Logf("value match: %v == %v", a, b)
		t.Fail()
	}
}

func Equal[T comparable](t *testing.T, a T, b T) {

	if a != b {
		t.Logf("value mismatch: %v != %v", a, b)
		t.Fail()
	}
}

func SlicesEqual[T comparable](t *testing.T, a []T, b []T) {

	if len(a) != len(b) {
		t.Logf("length mismatch: %v != %v", len(a), len(b))
		t.Fail()
	}

	for i := range a {
		if a[i] != b[i] {
			t.Logf("value mismatch at index %d: %v, %v", i, a[i], b[i])
			t.Fail()
			break
		}
	}
}

func NoErr(t *testing.T, err error) {

	if err != nil {
		t.Logf("error: %v", err)
		t.Fail()
	}
}
