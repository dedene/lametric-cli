package output

import (
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

// Spinner displays an animated spinner with a message.
type Spinner struct {
	frames  []string
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

// NewSpinner creates a new spinner with the given message.
func NewSpinner(msg string) *Spinner {
	return &Spinner{
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		message: msg,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation.
func (s *Spinner) Start() {
	// Only show spinner if stdout is a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		close(s.done)
		return
	}

	go func() {
		defer close(s.done)
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-s.stop:
				s.clear()
				return
			case <-ticker.C:
				s.mu.Lock()
				fmt.Printf("\r%s %s", s.frames[i%len(s.frames)], s.message)
				s.mu.Unlock()
				i++
			}
		}
	}()
}

// Stop stops the spinner and clears the line.
func (s *Spinner) Stop() {
	close(s.stop)
	<-s.done
}

// clear erases the spinner line.
func (s *Spinner) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Move to start of line and clear it
	fmt.Print("\r\033[K")
}
