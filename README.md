# timer [![GoDoc](https://godoc.org/github.com/jcuvillier/timer?status.svg)](https://godoc.org/github.com/jcuvillier/timer)

Timer a utility package to time function for performance tests.

```
go get github.com/jcuvillier/timer
```

Timer can execute a given function (`func() error`) a given number of times with parallelism.  
Executions are timed and a `Report` is returned.  

This report can be used to get statistics such as *min*, *max*, *mean* and *quantiles*.  
Also report can be printed out in a given `io.writer`. (See example below)

## Example

### Regular case

```golang
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
	r, err := timer.Run(f, 100, 5)
	if err != nil {
		log.Fatal(err)
	}
	// Print report to stdout
	r.Print(os.Stdout)
}

```
Output
```
Executions  [total, durations, rate]      100 10.4027215s 9.61
Durations   [min, mean, 50, 90, 99, max]  100.0798ms 511.142199ms 500.1322ms, 1.0001331s, 1.0002159s, 1.0002159s
Success     [ratio]                       88.00 %
Errors      [count]                       12
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
```
### With timeout

The `context.Context` can be used to stop the execution by cancelling the context, for instance with a timeout.

```golang
// Execute the function indefinitely with a parallelism of 5
// A context with timeout will interrupt the execution after 5 seconds
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
r, err := timer.Run(ctx, f, 0, 5)
if err != nil {
	log.Fatal(err)
}
// Print report to stdout
r.Print(os.Stdout)
```
Output
```
Executions  [total, durations, rate]      47 5.0001797s 9.40
Durations   [min, mean, 50, 90, 99, max]  100.0572ms 491.611323ms 400.1454ms, 900.1604ms, 1.0001512s, 1.0001512s
Success     [ratio]                       85.11 %
Errors      [count]                       7
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
   looks like this is an error
```


## Licence

MIT License

