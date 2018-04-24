package datatypes

type InPipes struct {
    pipes map[PipeId]*Event
}

// aux methods

func (inpipes InPipes) rinse(pipeid PipeId) {
    inpipes.pipes[pipeid] = nil
}

func (inpipes InPipes) consume(pipeid PipeId) (*Event, bool) {
    defer inpipes.rinse(pipeid)
    ev,ok := inpipes.pipes[pipeid]
    return ev,ok
}
