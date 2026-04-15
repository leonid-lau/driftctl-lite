package cmd

import (
	"bytes"
	"testing"
)

func executeCommand(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	// Reset flags between runs
	statePath = "terraform.tfstate"
	resource = ""
	outputFmt = "text"
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestRootCmd_DefaultFlags(t *testing.T) {
	out, err := executeCommand("--state", "terraform.tfstate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestRootCmd_ResourceFilter(t *testing.T) {
	out, err := executeCommand("--state", "terraform.tfstate", "--resource", "aws_s3_bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(out, "aws_s3_bucket") {
		t.Errorf("expected output to mention resource type, got: %s", out)
	}
}

func TestRootCmd_InvalidOutputFormat(t *testing.T) {
	_, err := executeCommand("--state", "terraform.tfstate", "--output", "xml")
	if err == nil {
		t.Fatal("expected error for unsupported output format, got nil")
	}
}

func TestRootCmd_JSONOutput(t *testing.T) {
	out, err := executeCommand("--state", "terraform.tfstate", "--output", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output for json format")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
