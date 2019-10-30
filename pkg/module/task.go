package module

// PipelineTask represents a task of a pipeline
type PipelineTask struct {
	Name string
	Image string
	Script string
	Triggers []PipelineTrigger
	Env []EnvEntry
}
