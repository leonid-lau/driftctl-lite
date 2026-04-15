package tfstate

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the state file watcher.
type WatchOptions struct {
	PollInterval time.Duration
}

// DefaultWatchOptions returns sensible defaults for WatchOptions.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		PollInterval: 5 * time.Second,
	}
}

// ChangeEvent is emitted when a watched state file changes.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a state file for changes and emits ChangeEvents.
type Watcher struct {
	path    string
	opts    WatchOptions
	lastSum string
}

// NewWatcher creates a Watcher for the given state file path.
func NewWatcher(path string, opts WatchOptions) *Watcher {
	return &Watcher{path: path, opts: opts}
}

// Watch starts polling and sends ChangeEvents on the returned channel.
// It stops when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) (<-chan ChangeEvent, error) {
	sum, err := hashFile(w.path)
	if err != nil {
		return nil, fmt.Errorf("watcher: initial hash failed: %w", err)
	}
	w.lastSum = sum

	ch := make(chan ChangeEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.opts.PollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				newSum, err := hashFile(w.path)
				if err != nil {
					continue
				}
				if newSum != w.lastSum {
					ch <- ChangeEvent{Path: w.path, OldHash: w.lastSum, NewHash: newSum}
					w.lastSum = newSum
				}
			}
		}
	}()
	return ch, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
