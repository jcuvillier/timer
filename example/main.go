package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jcuvillier/timer"
)

func main() {
	t := timer.New(100, 5)
	t.Run(func() error {
		r := rand.Int63n(10)
		time.Sleep(100*time.Millisecond + time.Duration(r)*100*time.Millisecond)
		if r == 0 {
			return fmt.Errorf("looks like this is an error")
		}
		return nil
	})
	t.Print(os.Stdout)
}
