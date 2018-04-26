package main

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
    "time"
)

func main() {

    inStreams, outStreams := changePointsExample()
    // inStreams, outStreams := clockExample()

    // Initialization
    inpipes := new(dt.InPipes)
    inpipes.Reset()
    var lastT dt.Time = -1 // minus infty

    for true {
        var nextT *dt.Time = nil
        // vote instreams
        for _, instr := range inStreams {
            aux := instr.StreamDef.PeekNextTime()
            nextT = dt.Min(aux, nextT)
        }
        // vote outstreams
        for _, outstr := range outStreams {
            aux := outstr.TicksDef.Vote(lastT)
            nextT = dt.Min(aux, nextT)
        }
        // end of execution
        if nextT == nil {
            break
        }
        // exec on input streams
        for _, instr := range inStreams {
            payload := instr.StreamDef.Exec(*nextT)
            inpipes.Put(instr.Name, dt.Event{*nextT, payload})
        }
        // exec on output streams
        for _, outstr:= range outStreams {
            payload := outstr.TicksDef.Exec(*nextT, *inpipes)
            if payload.IsSet {
                outpayload := outstr.ValDef.Exec(*nextT, payload.Val, *inpipes)
                inpipes.Put(outstr.Name, dt.Event{*nextT, outpayload})
            } else {
                inpipes.Put(outstr.Name, dt.Event{*nextT, dt.NothingPayload})
            }
        }
        // rinse output streams
        for _, outstr:= range outStreams {
            outstr.TicksDef.Rinse(*inpipes)
            outstr.ValDef.Rinse(*inpipes)
        }
        // reset pipes
        inpipes.Reset()
        lastT = *nextT
        time.Sleep(1000 * time.Millisecond)
    }

    fmt.Println("End of execution")
}
