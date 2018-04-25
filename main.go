package main

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
    "time"
)

type regularTicker struct {
    interval dt.Time;
    lastT dt.Time
}

func (ticker *regularTicker) PeekNextTime () dt.Time {
    return (ticker.lastT/ticker.interval)*ticker.interval+ticker.interval
}

func (ticker *regularTicker) Exec (t dt.Time) dt.EvPayload {
    ticker.lastT = t
    if t%ticker.interval == 0 {
        return dt.Some(dt.EpsVal{t/ticker.interval, nil})
    }
    return dt.NothingPayload
}

func main() {
    // Test streams
    oddsInStream := dt.InStream {"odds", &regularTicker{2,0}}
    allInStream := dt.InStream {"threes", &regularTicker{3,0}}
    inStreams := []dt.InStream{oddsInStream, allInStream}

    valodds := dt.OutStream{"oddsvals", dt.SrcTickerNode{"threes"}, &dt.PrevEqValNode{dt.TNode{}, "odds", []dt.Event{}}}
    shiftedvalodds := dt.OutStream{"shiftedthreesvals", &dt.DelayTickerNode{"threes", func (a dt.EvPayload, b dt.EvPayload)dt.EvPayload{return a}, []dt.Event{}}, &dt.PrevEqValNode{dt.TNode{}, "odds", *new([]dt.Event)}}
    outStreams := []dt.OutStream{valodds, shiftedvalodds}
    // Endof test streams

    inpipes := new(dt.InPipes)
    inpipes.Reset()
    var lastT dt.Time = -1 // minus infty

    for true {
        var nextT *dt.Time = nil
        // vote instreams
        for _, instr := range inStreams {
            aux := instr.StreamDef.PeekNextTime()
            nextT = dt.Min(&aux, nextT)
        }
        // vote outstreams
        for _, outstr := range outStreams {
            aux := outstr.TicksDef.Vote(lastT)
            nextT = dt.Min(aux, nextT)
        }
        // end of execution
        if nextT == nil {
            break // ???
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

    fmt.Println(oddsInStream.StreamDef.Exec(6))
}
