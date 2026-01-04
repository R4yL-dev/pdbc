package spinner

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a loading animation
type Spinner struct {
	message  string
	active   bool
	mu       sync.Mutex
	done     chan bool
	finished sync.WaitGroup
}

// New creates a new Spinner with the given message
func New(message string) *Spinner {
	return &Spinner{
		message: message,
		active:  false,
		done:    make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.finished.Add(1)
	s.mu.Unlock()

	go func() {
		defer s.finished.Done()
		frames := []string{".", "..", "..."}
		i := 0

		for {
			select {
			case <-s.done:
				return
			default:
				// Print message with current frame (overwriting previous line)
				// Use ANSI escape codes: \033[2K clears line, \r returns to start
				fmt.Printf("\033[2K\r%s%s", s.message, frames[i])
				i = (i + 1) % len(frames)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner animation and clears the line
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	// Signal goroutine to stop
	s.done <- true

	// Wait for goroutine to finish
	s.finished.Wait()

	// Clear the spinner line completely using ANSI escape codes
	// \033[2K clears the entire line, \r returns cursor to start
	fmt.Print("\033[2K\r")
}
