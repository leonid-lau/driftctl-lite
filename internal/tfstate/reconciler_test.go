package tfstate

import (
	"os"
	"testing"
	"time"
)

func newTestReconciler(t *testing.T, opts ReconcileOptions) (*Reconciler, *SnapshotStore) {
	t.Helper()
	dir, err := os.MkdirTemp("", "reconciler-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	store, err := NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	cache := NewStateCache(opts.MaxAge)
	return NewReconciler(store, cache, opts), store
}

func TestReconcile_NoPreviousSnapshot(t *testing.T) {
	r, _ := newTestReconciler(t, DefaultReconcileOptions())

	state := &State{Resources: []Resource{
		{Type: "aws_s3_bucket", Name: "my-bucket", Attributes: map[string]interface{}{"region": "us-east-1"}},
	}}

	res, err := r.Reconcile("prod", state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FromCache {
		t.Error("expected FromCache=false on first run")
	}
	if len(res.Diffs) != 0 {
		t.Errorf("expected no diffs on first run, got %d", len(res.Diffs))
	}
}

func TestReconcile_NoDrift(t *testing.T) {
	r, _ := newTestReconciler(t, DefaultReconcileOptions())

	state := &State{Resources: []Resource{
		{Type: "aws_s3_bucket", Name: "b", Attributes: map[string]interface{}{"acl": "private"}},
	}}

	// First run — seeds the store.
	if _, err := r.Reconcile("env", state); err != nil {
		t.Fatalf("first reconcile: %v", err)
	}

	// Second run — identical state should return FromCache=true.
	res, err := r.Reconcile("env", state)
	if err != nil {
		t.Fatalf("second reconcile: %v", err)
	}
	if !res.FromCache {
		t.Error("expected FromCache=true when state is unchanged")
	}
}

func TestReconcile_DetectsDrift(t *testing.T) {
	r, _ := newTestReconciler(t, DefaultReconcileOptions())

	original := &State{Resources: []Resource{
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"instance_type": "t2.micro"}},
	}}
	if _, err := r.Reconcile("staging", original); err != nil {
		t.Fatalf("seed reconcile: %v", err)
	}

	drifted := &State{Resources: []Resource{
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"instance_type": "t2.large"}},
	}}
	res, err := r.Reconcile("staging", drifted)
	if err != nil {
		t.Fatalf("drift reconcile: %v", err)
	}
	if res.FromCache {
		t.Error("expected FromCache=false when state changed")
	}
	if len(res.Diffs) == 0 {
		t.Error("expected diffs to be reported")
	}
}

func TestReconcile_DryRun_DoesNotPersist(t *testing.T) {
	opts := ReconcileOptions{DryRun: true, MaxAge: time.Minute}
	r, store := newTestReconciler(t, opts)

	state := &State{Resources: []Resource{
		{Type: "aws_vpc", Name: "main", Attributes: map[string]interface{}{"cidr": "10.0.0.0/16"}},
	}}
	if _, err := r.Reconcile("vpc-key", state); err != nil {
		t.Fatalf("dry-run reconcile: %v", err)
	}

	// Store should have nothing persisted.
	if _, err := store.Load("vpc-key"); err == nil {
		t.Error("expected load to fail after dry-run, but snapshot was persisted")
	}
}

func TestReconcile_NilState_ReturnsError(t *testing.T) {
	r, _ := newTestReconciler(t, DefaultReconcileOptions())
	if _, err := r.Reconcile("k", nil); err == nil {
		t.Error("expected error for nil state")
	}
}
