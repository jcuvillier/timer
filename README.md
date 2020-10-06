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
```
Prints:
```
Executions   [total, durations, rate]       100 10.4025927s 9.61
Durations    [min, mean, 50, 90, 99, max]   100.0695ms 511.130443ms 500.1308ms, 1.0000846s, 1.0002318s, 1.0002318s
Success      [ratio]                        88.00 %
Errors       [count]                        12
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

