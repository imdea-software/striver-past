package datatypes

type InPipes struct {
    pipes map[StreamName]Event
}

// aux methods

func (inpipes InPipes) Reset() {
    inpipes.pipes = make(map[StreamName]Event)
}

func (inpipes InPipes) strictConsume(streamId StreamName) Event {
    ev,ok := inpipes.pipes[streamId]
    if !ok {
        panic("Failed strict consume")
    }
    return ev
}
