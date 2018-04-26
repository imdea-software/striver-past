package main

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
    "time"
)

func main() {

    inStreams, outStreams := shiftExample()
    // inStreams, outStreams := changePointsExample()
    //inStreams, outStreams := clockExample()

    // Initialization
    inpipes := new(dt.InPipes)
    inpipes.Reset()
    var lastT dt.Time = -1 // minus infty

    for true {
        var nextT dt.MaybeTime = dt.NothingTime
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
        if !nextT.IsSet {
            break
        }
        // exec on input streams
        for _, instr := range inStreams {
            payload := instr.StreamDef.Exec(nextT.Val)
            inpipes.Put(instr.Name, dt.Event{nextT.Val, payload})
        }
        // exec on output streams
        for _, outstr:= range outStreams {
            payload := outstr.TicksDef.Exec(nextT.Val, *inpipes)
            if payload.IsSet {
                outpayload := outstr.ValDef.Exec(nextT.Val, payload.Val, *inpipes)
                inpipes.Put(outstr.Name, dt.Event{nextT.Val, outpayload})
            } else {
                inpipes.Put(outstr.Name, dt.Event{nextT.Val, dt.NothingPayload})
            }
        }
        // rinse output streams
        for _, outstr:= range outStreams {
            outstr.TicksDef.Rinse(*inpipes)
            outstr.ValDef.Rinse(*inpipes)
        }
        // reset pipes
        inpipes.Reset()
        lastT = nextT.Val
        time.Sleep(1000 * time.Millisecond)
    }

    fmt.Println("End of execution")
}
