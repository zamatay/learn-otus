package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	fn := func(out Out, bi Bi) {
		defer close(bi)

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
	}

	for _, stage := range stages {
		twin := make(Bi)

		go fn(in, twin)

		in = stage(twin)
	}

	return in
}
