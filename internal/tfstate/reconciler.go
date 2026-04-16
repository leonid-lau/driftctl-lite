package tfstate

import (
	"fmt"
	"time"
)

// ReconcileOptions controls reconciliation behaviour.
type ReconcileOptions struct {
	// DryRun skips persisting the reconciled snapshot.
	DryRun bool
	// MaxAge is the maximum age of a cached snapshot before it is considered stale.
	MaxAge time.Duration
}

// DefaultReconcileOptions returns sensible defaults.
func DefaultReconcileOptions() ReconcileOptions {
	return ReconcileOptions{
		DryRun: false,
		MaxAge: 5 * time.Minute,
	}
}

// ReconcileResult holds the outcome of a reconciliation run.
type ReconcileResult struct {
	Snapshot  *Snapshot
	FromCache bool
	Diffs     []ResourceDiff
}

// Reconciler compares a freshly loaded state against the last known
// snapshot and records the new snapshot when changes are detected.
type Reconciler struct {
	store  *SnapshotStore
	cache  *StateCache
	opts   ReconcileOptions
}

// NewReconciler constructs a Reconciler backed by the given store and cache.
func NewReconciler(store *SnapshotStore, cache *StateCache, opts ReconcileOptions) *Reconciler {
	return &Reconciler{store: store, cache: cache, opts: opts}
}

// Reconcile loads the current state for key, diffs it against the previous
// snapshot, persists a new snapshot (unless DryRun), and returns the result.
func (r *Reconciler) Reconcile(key string, current *State) (*ReconcileResult, error) {
	if current == nil {
		return nil, fmt.Errorf("reconciler: current state must not be nil")
	}

	newSnap := NewSnapshot(current.Resources)

	prev, err := r.store.Load(key)
	if err != nil {
		// No previous snapshot — treat all resources as new (no diffs).
		prev = nil
	}

	result := &ReconcileResult{Snapshot: newSnap}

	if prev != nil && prev.Equal(newSnap) {
		result.FromCache = true
		return result, nil
	}

	// Compute attribute-level diffs for changed resources.
	if prev != nil {
		for id, res := range newSnap.Resources {
			if old, ok := prev.Resources[id]; ok {
				diffs := DiffAttributes(old.Attributes, res.Attributes)
				result.Diffs = append(result.Diffs, diffs...)
			}
		}
	}

	if !r.opts.DryRun {
		if err := r.store.Save(key, newSnap); err != nil {
			return nil, fmt.Errorf("reconciler: save snapshot: %w", err)
		}
	}

	return result, nil
}
