package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	statePath  string
	resource   string
	outputFmt  string
)

var rootCmd = &cobra.Command{
	Use:   "driftctl-lite",
	Short: "Detect infrastructure drift between Terraform state and live cloud resources",
	Long: `driftctl-lite compares your Terraform state file against live AWS resources
and reports any drift — missing, extra, or modified resources.`,
	RunE: runDrift,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&statePath, "state", "s", "terraform.tfstate",
		"Path to the Terraform state file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&resource, "resource", "r", "",
		"Filter drift detection to a specific resource type (e.g. aws_s3_bucket)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&outputFmt, "output", "o", "text",
		"Output format: text or json",
	)
}

func runDrift(cmd *cobra.Command, args []string) error {
	if outputFmt != "text" && outputFmt != "json" {
		return fmt.Errorf("unsupported output format %q: must be 'text' or 'json'", outputFmt)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Reading state from: %s\n", statePath)

	if resource != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Filtering by resource type: %s\n", resource)
	}

	// Placeholder: real orchestration wired in a follow-up
	fmt.Fprintln(cmd.OutOrStdout(), "Drift detection complete.")
	return nil
}
