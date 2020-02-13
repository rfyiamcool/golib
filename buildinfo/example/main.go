package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rfyiamcool/golib/buildinfo"
)

func main() {
	v := flag.Bool("v", false, "show bin info")
	flag.Parse()
	if *v {
		_, _ = fmt.Fprint(os.Stderr, buildinfo.StringifyMultiLine())
		os.Exit(1)
	}

	fmt.Println("my app running...")
}
