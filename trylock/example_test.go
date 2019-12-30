package trylock_test

import (
	"fmt"

	"github.com/rfyiamcool/golib/trylock"
)

func Example() {
	var mu trylock.Mutex
	fmt.Println(mu.TryLock())
	fmt.Println(mu.TryLock())
	mu.Unlock()
	fmt.Println(mu.TryLock())
	// Output:
	// true
	// false
	// true
}
