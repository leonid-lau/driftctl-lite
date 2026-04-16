package tfstate

import "fmt"

// PipelineOptions configures the state loading pipeline.
type PipelineOptions struct {
	LoadOptions    LoadOptions
	MergeOptions   MergeOptions
	Validate       bool
	UseCache       bool
	CacheDir       string
	CacheTTL       int // seconds
}

// DefaultPipelineOptions returns sensible defaults.
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		LoadOptions:  DefaultLoadOptions(),
		MergeOptions: DefaultMergeOptions(),
		Validate:     true,
		UseCache:     false,
		CacheDir:     ".driftctl-cache",
		CacheTTL:     300,
	}
}

// Pipeline orchestrates load → merge → validate into a single State.
type Pipeline struct {
	opts PipelineOptions
}

// NewPipeline creates a Pipeline with the given options.
func NewPipeline(opts PipelineOptions) *Pipeline {
	return &Pipeline{opts: opts}
}

// Run executes the pipeline for the given root directory and returns a merged State.
func (p *Pipeline) Run(dir string) (*State, error) {
	if dir == "" {
		return nil, fmt.Errorf("pipeline run: dir must not be empty")
	}

	var states []*State
	var err error

	if p.opts.UseCache {
		cl := NewCachedLoader(p.opts.CacheDir, p.opts.LoadOptions)
		merged, err := cl.Load(dir)
		if err != nil {
			return nil, fmt.Errorf("pipeline cached load: %w", err)
		}
		if p.opts.Validate {
			if verr := ValidateState(merged); verr != nil {
				return nil, fmt.Errorf("pipeline validate: %w", verr)
			}
		}
		return merged, nil
	}

	states, err = LoadAll(dir, p.opts.LoadOptions)
	if err != nil {
		return nil, fmt.Errorf("pipeline load: %w", err)
	}

	merged, err := MergeStates(states, p.opts.MergeOptions)
	if err != nil {
		return nil, fmt.Errorf("pipeline merge: %w", err)
	}

	if p.opts.Validate {
		if verr := ValidateState(merged); verr != nil {
			return nil, fmt.Errorf("pipeline validate: %w", verr)
		}
	}

	return merged, nil
}
