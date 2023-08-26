package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		twin := make(Bi)

		go func(bi Bi, out Out) {
			defer close(twin)

			for {
				select {
				case value, isOpen := <-out:
					if !isOpen {
						return
					}
					bi <- value
				case <-done:
					return
				}
			}
		}(twin, out)

		out = stage(twin)
	}

	return out
}
