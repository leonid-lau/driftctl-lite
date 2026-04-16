package tfstate_test

import (
	"os"
	"testing"

	"github.com/example/driftctl-lite/internal/tfstate"
)

// TestReconciler_FullCycle exercises the reconciler end-to-end: seed →
// unchanged → drifted, verifying snapshot persistence across instances.
func TestReconciler_FullCycle(t *testing.T) {
	dir, err := os.MkdirTemp("", "reconciler-integration-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	store, err := tfstate.NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	opts := tfstate.DefaultReconcileOptions()
	r := tfstate.NewReconciler(store, tfstate.NewStateCache(opts.MaxAge), opts)

	base := &tfstate.State{
		Resources: []tfstate.Resource{
			{Type: "aws_s3_bucket", Name: "logs", Attributes: map[string]interface{}{"versioning": "enabled"}},
			{Type: "aws_iam_role", Name: "deployer", Attributes: map[string]interface{}{"path": "/"}},
		},
	}

	// --- Seed ---
	res1, err := r.Reconcile("prod", base)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	if res1.FromCache {
		t.Error("seed: expected FromCache=false")
	}

	// --- No change ---
	res2, err := r.Reconcile("prod", base)
	if err != nil {
		t.Fatalf("no-change: %v", err)
	}
	if !res2.FromCache {
		t.Error("no-change: expected FromCache=true")
	}

	// --- Drift ---
	drifted := &tfstate.State{
		Resources: []tfstate.Resource{
			{Type: "aws_s3_bucket", Name: "logs", Attributes: map[string]interface{}{"versioning": "suspended"}},
			{Type: "aws_iam_role", Name: "deployer", Attributes: map[string]interface{}{"path": "/"}},
		},
	}
	res3, err := r.Reconcile("prod", drifted)
	if err != nil {
		t.Fatalf("drift: %v", err)
	}
	if res3.FromCache {
		t.Error("drift: expected FromCache=false")
	}
	if len(res3.Diffs) == 0 {
		t.Error("drift: expected at least one attribute diff")
	}

	// Verify new snapshot was persisted by loading from a fresh store instance.
	store2, _ := tfstate.NewSnapshotStore(dir)
	snap, err := store2.Load("prod")
	if err != nil {
		t.Fatalf("load after drift: %v", err)
	}
	if snap == nil {
		t.Fatal("expected persisted snapshot, got nil")
	}
}
