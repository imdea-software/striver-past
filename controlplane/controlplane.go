package controlplane

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    //"time"
)

func Start(inStreams []dt.InStream, outStreams []dt.OutStream, outchan chan dt.FlowingEvent) {

    // Initialization
    inpipes := dt.InPipes{make(map[dt.StreamName]dt.Event), outchan}
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
            payload := outstr.TicksDef.Exec(nextT.Val, inpipes)
            if payload.IsSet {
                outpayload := outstr.ValDef.Exec(nextT.Val, payload.Val, inpipes)
                inpipes.Put(outstr.Name, dt.Event{nextT.Val, outpayload})
            } else {
                inpipes.Put(outstr.Name, dt.Event{nextT.Val, dt.NothingPayload})
            }
        }
        // rinse output streams
        for _, outstr:= range outStreams {
            outstr.TicksDef.Rinse(inpipes)
            outstr.ValDef.Rinse(inpipes)
        }
        // reset pipes
        inpipes.Reset()
        lastT = nextT.Val
        //time.Sleep(1000 * time.Millisecond)
    }
}
