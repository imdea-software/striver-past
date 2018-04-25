package datatypes

import "fmt"

type InPipes struct {
    Pipes map[StreamName]Event
}

// aux methods

func (inpipes *InPipes) Reset() {
    inpipes.Pipes = make(map[StreamName]Event)
}

func (inpipes InPipes) strictConsume(streamId StreamName) Event {
    ev,ok := inpipes.Pipes[streamId]
    if !ok {
        panic("Failed strict consume")
    }
    return ev
}

func (inpipes InPipes) Put(streamId StreamName, ev Event) {
    fmt.Printf("Stream %s saying %v\n", streamId, ev)
    inpipes.Pipes[streamId] = ev
}
