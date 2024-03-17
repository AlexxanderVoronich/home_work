package hw06pipelineexecution

type (
	// In represents the input channel for each stage of the pipeline.
	In = <-chan interface{}
	// Out represents the output channel for each stage of the pipeline.
	Out = In
	// Bi represents the bidirectional channel for passing data between stages.
	Bi = chan interface{}
)

// Stage represents a separate stage in the pipelines.
type Stage func(in In) (out Out)

// ExecutePipeline links the given stages in pipeline.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageWrapper := func(in In, done In, stage Stage) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			outInternal := stage(in)

			for {
				select {
				case v, ok := <-outInternal:
					if !ok {
						return
					}
					if v != nil {
						out <- v
					}
				case <-done:
					return
				}
			}
		}()
		return out
	}

	out := in
	for _, stage := range stages {
		out = stageWrapper(out, done, stage)
	}
	return out
}
