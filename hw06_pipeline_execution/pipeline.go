package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageIn := make(Bi)
		go func(in In) {
			defer close(stageIn)
			for {
				select {
				case <-done:
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					stageIn <- v
				}
			}
		}(in)

		in = stage(stageIn)
	}

	return in
}
