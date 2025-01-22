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
		if stage != nil {
			out = stageBuilder(out, done, stage)
		}
	}

	return out
}

func stageBuilder(in In, done In, stage Stage) Out {
	out := make(Bi)

	go func() {
		for {
			select {
			case <-done:
				go func() {
					for range in {
					}
				}()
				close(out)
				return
			case v, ok := <-in:
				if !ok {
					close(out)
					return
				}
				out <- v
			}
		}
	}()

	return stage(out)
}
