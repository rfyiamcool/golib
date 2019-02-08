package tester

import (
	"fmt"
	"testing"
)

func Example_caller() {
	f := func() {
		file, line := caller()
		fmt.Printf("%s:%d", file, line)
	}
	f()
}

type stTest struct{}
type stTestInterface interface{}

func TestExpectReject(t *testing.T) {
	// Standard expectations
	Expect(t, "a", "a")
	Expect(t, 42, 42)
	Expect(t, nil, nil)
	Expect(t, stTestInterface(nil), nil)
	Expect(t, []int{42}, []int{42})

	// Standard rejections
	Reject(t, "a", "A")
	Reject(t, 42, int64(42.0))
	Reject(t, 42, 42.0)
	Reject(t, 42, "42")
	Reject(t, []int{42}, []int{41})
	Reject(t, stTest{}, nil)
	Reject(t, []string{}, nil)
	Reject(t, []stTest{}, nil)

	var typedNil *stTest
	Reject(t, typedNil, nil)

	// Table-based test
	examples := []struct{ a, b string }{
		{"first", "first"},
		{"second", "second"},
	}

	for i, ex := range examples {
		Expect(t, ex, ex, i)
		Expect(t, &ex, &ex, i)

		Reject(t, ex, &ex, i)
		Reject(t, ex, 0, i)
		Reject(t, ex, "", i)
		Reject(t, ex, byte('a'), i)
		Reject(t, ex, float64(5.9), i)
	}
}

func TestAssertRefute(t *testing.T) {
	// Standard assertions
	Assert(t, "a", "a")
	Assert(t, 42, 42)
	Assert(t, nil, nil)
	Assert(t, []int{42}, []int{42})

	// Standard refutations
	Refute(t, "a", "A")
	Refute(t, 42, int64(42.0))
	Refute(t, 42, 42.0)
	Refute(t, 42, "42")
	Refute(t, []int{42}, []int{41})
	Refute(t, []string{}, nil)
	Refute(t, []stTest{}, nil)

	// Table-based test
	examples := []struct{ a, b string }{
		{"first", "first"},
		{"second", "second"},
	}

	// Note: there's no argument to pass the index to assertions.
	for _, ex := range examples {
		Assert(t, ex, ex)
		Assert(t, &ex, &ex)

		Refute(t, ex, &ex)
		Refute(t, ex, 0)
		Refute(t, ex, "")
		Refute(t, ex, byte('a'))
		Refute(t, ex, float64(5.9))
	}
}

func Test_exampleNum(t *testing.T) {
	expectationFunc := func(t *testing.T, n ...int) []int {
		return n
	}

	Expect(t, exampleNum(expectationFunc(t)), "")
	Expect(t, exampleNum(expectationFunc(t, 0)), "0.")
	Expect(t, exampleNum(expectationFunc(t, 1)), "1.")
	Expect(t, exampleNum(expectationFunc(t, 2)), "2.")
}
