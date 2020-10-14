package timer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"
)

// Report holds execution results.
type Report struct {
	num       int
	durations []time.Duration
	errs      []error
	total     time.Duration
}

// Run executes the given function *n* times concurently with a parallelism of *p*.
// The executions can also be stopped by cancelling the context.
//
// A negative or zero *n* means infitnite loop.
// Note that this mode is meant to be used with a cancellable context.
//
// *p* must be positive
//
// Each execution is timed and potential errors saved in the returned report.
// Total execution time is also saved in the report.
func Run(ctx context.Context, f func() error, n, p int) (*Report, error) {
	// Check arguments
	if p <= 0 {
		return nil, errors.New("p must be positive")
	}

	// Prepare execution and report
	r := &Report{
		durations: []time.Duration{},
		errs:      []error{},
	}
	poolC := make(chan bool, p)      // Buffered channel used as a pool to control the concurrency
	durC := make(chan time.Duration) // Channel to get durations from execution go routines
	errC := make(chan error)         // Channel to get errors from execution go routines
	doneC := make(chan bool)         // Channel to control end of execution
	k := 0

	// Start the executions
	startT := time.Now()
	go func() {
		i := 0
		for {
			if n > 0 && i > n {
				break
			}
			poolC <- true
			go execAndTimeFunction(f, durC, errC, poolC, doneC)
		}
	}()

	// Read channels and save results
L:
	for {
		select {
		case d := <-durC:
			r.durations = append(r.durations, d)
		case err := <-errC:
			r.errs = append(r.errs, err)
		case <-doneC:
			k = k + 1
			if n > 0 && k >= n {
				break L
			}
		case <-ctx.Done():
			break L
		}
	}
	r.total = time.Since(startT)
	r.num = k

	// Sort r.durations before returning
	sort.Slice(r.durations, func(i, j int) bool { return r.durations[i] < r.durations[j] })
	return r, nil
}

// execAndTimeFunction executes the given function
//
// the execution is timed and its duration is sent into the durC channel
// potential error is sent to errC channel
func execAndTimeFunction(f func() error, durC chan time.Duration, errC chan error, poolC chan bool, doneC chan bool) {
	// Execute f and get send its duration to durC
	startT := time.Now()
	err := f()
	durC <- time.Since(startT)

	// Send error to errC if applicable
	if err != nil {
		errC <- err
	}

	// Handle pool and done channel
	<-poolC
	doneC <- true
}

// Executions returns the total number of execution done
func (r *Report) Executions() int {
	return r.num
}

// Duration returns the total duration of the executions
func (r *Report) Duration() time.Duration {
	return r.total
}

// Rate returns the mean number of execution per second
func (r *Report) Rate() float64 {
	return float64(r.num) * float64(time.Second) / float64(r.total)
}

// Min returns the minimum duration value
func (r *Report) Min() time.Duration {
	return r.durations[0]
}

// Max returns the maximum duration value
func (r *Report) Max() time.Duration {
	return r.durations[len(r.durations)-1]
}

// Mean returns the mean duration value
func (r *Report) Mean() time.Duration {
	var total time.Duration
	for _, d := range r.durations {
		total = total + d
	}
	return time.Duration(int64(total) / int64(len(r.durations)))
}

// Quantile returns the given quantile value
func (r *Report) Quantile(q int) time.Duration {
	if q < 0 || q > 100 {
		panic(fmt.Sprintf("quantile must be between 0 and 100, given %d", q))
	}
	index := q * len(r.durations) / 100
	return r.durations[index]
}

// Errors returns thrown errors
func (r *Report) Errors() []error {
	return r.errs
}

// Print write a report in the given writer.
// Format is:
//
// Executions [total, duration, rate]
// Durations  [min, mean, 50, 90, 99, max]
// Success    [ratio]
// Errors     [count]
//
func (r *Report) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Executions\t[total, durations, rate]\t%d %s %.2f\n", r.Executions(), r.Duration(), r.Rate())
	fmt.Fprintf(tw, "Durations\t[min, mean, 50, 90, 99, max]\t%s %s %s, %s, %s, %s\n", r.Min(), r.Mean(), r.Quantile(50), r.Quantile(90), r.Quantile(99), r.Max())
	fmt.Fprintf(tw, "Success\t[ratio]\t%.2f %%\n", float64(100)-float64(len(r.errs)*100)/float64(r.num))
	fmt.Fprintf(tw, "Errors\t[count]\t%d\n", len(r.errs))
	tw.Flush()
	for _, err := range r.errs {
		fmt.Fprintf(w, "   %s\n", err.Error())
	}
}
