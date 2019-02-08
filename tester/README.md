## tester

A tiny test framework for making short, useful assertions in your Go tests.

Assert(t, have, want) and Refute(t, have, want) abort a test immediately with t.Fatal.

Expect(t, have, want) and Reject(t, have, want) allow a test to continue, reporting failure at the end with t.Error.

They print nice error messages, preserving the order of have (actual result) before want (expected result) to minimize confusion.

## Usage

```
func TestExample(t *testing.T) {
	tester.Expect(t, "a", "a")
	tester.Reject(t, 42, int64(42))

	tester.Assert(t, "b", "b")
	tester.Refute(t, 99, int64(99))
}

func TestTableExample(t *testing.T) {
	examples := []struct{ a, b string }{
		{"first", "first"},
		{"second", "second"},
	}

	// Pass the index to improve the error message for table-based tests.
	for i, ex := range examples {
		tester.Expect(t, ex, ex, i)
		tester.Reject(t, ex, &ex, i)
	}

	// Cannot pass index into Assert or Refute, they fail fast.
	for _, ex := range examples {
		tester.Assert(t, ex, ex)
		tester.Refute(t, ex, &ex)
	}
}
```
