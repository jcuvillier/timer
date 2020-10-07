# timer
Timer a utility package to time function for performance tests.

```
go get github.com/jcuvillier/timer
```

## Example

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
```
Prints:
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

