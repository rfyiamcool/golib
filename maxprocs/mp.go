package maxprocs

import (
	"go.uber.org/automaxprocs/maxprocs"
)

func AutoMaxProcess() {
	nopLog := func(string, ...interface{}) {}
	maxprocs.Set(maxprocs.Logger(nopLog))
}
