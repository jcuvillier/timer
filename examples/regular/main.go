package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jcuvillier/timer"
)

func main() {
	// Create function to test
	// This function will sleep 100ms + x*100ms, x being picked randomly between 0 and 10
	// If x equals 0, an error is raised to show how error are handled
	f := func() error {
		rnd := rand.Int63n(10)
		time.Sleep(100*time.Millisecond + time.Duration(rnd)*100*time.Millisecond)
		if rnd == 0 {
			return fmt.Errorf("looks like this is an error")
		}
		return nil
	}
	// Execute the function 100 times with a parallelism of 5
	r, err := timer.Run(context.Background(), f, -1, 5)
	if err != nil {
		log.Fatal(err)
	}
	// Print report to stdout
	r.Print(os.Stdout)
}
