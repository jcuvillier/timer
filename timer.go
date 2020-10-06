package timer

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/tabwriter"
	"time"
)

// Timer
type Timer struct {
	num  int
	conc int

	durations []time.Duration
	errs      []error
	wg        *sync.WaitGroup
	total     time.Duration
}

func New(n, c int) *Timer {
	return &Timer{
		num:       n,
		conc:      c,
		durations: make([]time.Duration, 0, n),
		errs:      make([]error, 0, n),
		wg:        &sync.WaitGroup{},
	}
}

func (t *Timer) Run(f func() error) {
	t.wg.Add(t.num)
	poolC := make(chan bool, t.conc)
	durC := make(chan time.Duration)
	errC := make(chan error)
	doneC := make(chan bool)

	startT := time.Now()

	go func() {
		for i := 0; i < t.num; i++ {
			poolC <- true
			go t.execAndTimeFunction(f, durC, errC, poolC)
		}
	}()

	// Wait for completion
	go func() {
		t.wg.Wait()
		doneC <- true
	}()

	// Read channels and save results
L:
	for {
		select {
		case d := <-durC:
			t.durations = append(t.durations, d)
		case err := <-errC:
			t.errs = append(t.errs, err)
		case <-doneC:
			break L
		}
	}
	t.total = time.Since(startT)

	// Sort t.durations before returning
	sort.Slice(t.durations, func(i, j int) bool { return t.durations[i] < t.durations[j] })
	return
}

// execAndTimeFunction executes the given function
// the execution is timed and its duration is sent into the durC channel
// potential error is sent to errC channel
func (t *Timer) execAndTimeFunction(f func() error, durC chan time.Duration, errC chan error, poolC chan bool) {
	// Execute f and get send its duration to durC
	startT := time.Now()
	err := f()
	durC <- time.Since(startT)

	// Send error to errC if applicable
	if err != nil {
		errC <- err
	}

	// Handle pool channel and wait group
	<-poolC
	t.wg.Done()
}

// Executions returns the total number of execution done
func (t *Timer) Executions() int {
	return t.num
}

// Duration returns the total duration of the executions
func (t *Timer) Duration() time.Duration {
	return t.total
}

// Rate returns the mean number of execution per second
func (t *Timer) Rate() float64 {
	return float64(t.num) * float64(time.Second) / float64(t.total)
}

// Min returns the minimum duration value
func (t *Timer) Min() time.Duration {
	return t.durations[0]
}

// Max returns the maximum duration value
func (t *Timer) Max() time.Duration {
	return t.durations[len(t.durations)-1]
}

// Mean returns the mean duration value
func (t *Timer) Mean() time.Duration {
	var total time.Duration
	for _, d := range t.durations {
		total = total + d
	}
	return time.Duration(int64(total) / int64(len(t.durations)))
}

// Quantile returns the given quantile value
func (t *Timer) Quantile(q int) time.Duration {
	if q < 0 || q > 100 {
		panic(fmt.Sprintf("quantile must be between 0 and 100, given %d", q))
	}
	index := q * len(t.durations) / 100
	return t.durations[index]
}

// Print write a report in the given writer
// Format is:
//
// Executions [total, duration, rate]      ~~ ~~ ~~
// Durations  [min, mean, 50, 90, 99, max] ~~ ~~ ~~ ~~ ~~ ~~
// Success    [ratio]                      ~~%
// Errors     [count]					   ~~
//   ~~~~~~~~~~~~~~~~~~~~~~
//   ~~~~~~~~~~~~~~~~~~~~~~
//   ~~~~~~~~~~~~~~~~~~~~~~
func (t *Timer) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintf(tw, "Executions\t[total, durations, rate]\t%d %s %.2f\n", t.Executions(), t.Duration(), t.Rate())
	fmt.Fprintf(tw, "Durations\t[min, mean, 50, 90, 99, max]\t%s %s %s, %s, %s, %s\n", t.Min(), t.Mean(), t.Quantile(50), t.Quantile(90), t.Quantile(99), t.Max())
	fmt.Fprintf(tw, "Success\t[ratio]\t%.2f %%\n", float64(100)-float64(len(t.errs)*100)/float64(t.num))
	fmt.Fprintf(tw, "Errors\t[count]\t%d\n", len(t.errs))
	tw.Flush()
	for _, err := range t.errs {
		fmt.Fprintf(w, "   %s\n", err.Error())
	}
}
