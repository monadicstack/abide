package wait

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// InterruptSignal creates a channel that will accept the next SIGTERM or SIGINT
// signal the OS sends to this process. This call does NOT block - it is up to you
// to "<-" in or out of a 'select' to actually control program flow.
func InterruptSignal() chan os.Signal {
	var interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return interrupt
}

// ContextOrGroupOrInterrupt blocks until one of 3 things happens: the context
// deadline is met (if one exists), the wait group's Wait() function naturally
// unblocks, or we receive a SIGINT/SIGTERM signal.
func ContextOrGroupOrInterrupt(ctx context.Context, wg *sync.WaitGroup) {
	waitGroupDone := make(chan bool, 1)
	unblock := make(chan bool, 1)

	go func() {
		// We're going to hold off on letting the outer function finish until one of
		// our 3 conditions have been met; context deadline, wait group finished, or interrupt.
		select {
		case <-ctx.Done():
			// We hit the deadline for our context.
		case <-waitGroupDone:
			// The WaitGroup finished naturally, so unblock this goroutine w/o worrying about the context.
		case <-InterruptSignal():
			// We received a second interrupt/terminate signal so be kind and give up. This
			// helps prevent your program from getting stuck in this function regardless of
			// how many times the user hits Ctrl+C.
		}
		unblock <- true
	}()

	go func() {
		wg.Wait()
		waitGroupDone <- true
	}()

	<-unblock
}
