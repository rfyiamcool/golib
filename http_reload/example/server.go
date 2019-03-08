// +build go1.8

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rfyiamcool/golib/http_reload"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := r.URL.Query().Get("duration")
		if len(d) != 0 {
			t, _ := time.ParseDuration(d)
			time.Sleep(t)
		}
		fmt.Fprintln(w, "Hello, World!")
	})

	log.Fatalln(httpReload.ListenAndServe(":8080", nil))
}
