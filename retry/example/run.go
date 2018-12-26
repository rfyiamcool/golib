package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/rfyiamcool/golib/retry"
)

func main() {
	r := retry.New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var running = false
	err := r.Ensure(ctx, func() error {
		log.Println("enter")
		if !running {
			log.Println("111")
			running = true
			return retry.Retriable(errors.New("diy"))
		}

		log.Println("222")
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
