package tfstate

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeTempWatchFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tfstate-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatcher_NoChangeEmitsNothing(t *testing.T) {
	path := writeTempWatchFile(t, `{"version":4}`)
	opts := WatchOptions{P * time.Millisecond(path, opts)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}
	<-ctx.Done()
	if len(ch) != 0 {
		t.Errorf("expected no events, got %d", len(ch))
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempWatchFile(t, `{"version":4}`)
	opts := WatchOptions{PollInterval: 20 * time.Millisecond}
	w := NewWatcher(path, opts)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Overwrite the file after a short delay.
	time.AfterFunc(60*time.Millisecond, func() {
		if err := os.WriteFile(path, []byte(`{"version":4,"serial":2}`), 0644); err != nil {
			t.Errorf("write file: %v", err)
		}
	})

	select {
	case evt := <-ch:
		if evt.Path != path {
			t.Errorf("expected path %s, got %s", path, evt.Path)
		}
		if evt.OldHash == evt.NewHash {
			t.Errorf("expected hashes to differ")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_InvalidPath(t *testing.T) {
	opts := DefaultWatchOptions()
	w := NewWatcher("/nonexistent/path/state.tfstate", opts)
	_, err := w.Watch(context.Background())
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestDefaultWatchOptions(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.PollInterval != 5*time.Second {
		t.Errorf("unexpected default poll interval: %v", opts.PollInterval)
	}
}
