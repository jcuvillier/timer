package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jcuvillier/timer"
)

func main() {
	f := func() error {
		rnd := rand.Int63n(10)
		time.Sleep(100*time.Millisecond + time.Duration(rnd)*100*time.Millisecond)
		if rnd == 0 {
			return fmt.Errorf("looks like this is an error")
		}
		return nil
	}
	r, err := timer.Run(f, 100, 5)
	if err != nil {
		log.Fatal(err)
	}
	r.Print(os.Stdout)
}
