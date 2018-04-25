package main

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "fmt"
)

type regularTicker struct {
    interval dt.Time;
    lastT dt.Time
}

func (ticker regularTicker) PeekNextTime () dt.Time { return (ticker.lastT/ticker.interval)*ticker.interval+ticker.interval }
func (ticker regularTicker) Exec (t dt.Time) dt.EvPayload {
    ticker.lastT = t
    if t%ticker.interval == 0 {
        return dt.Some(t/ticker.interval)
    }
    return dt.NothingPayload
}

func main() {

    oddsInStream := dt.InStream {"odds", regularTicker{2,0}}

    inpipes := *new(dt.InPipes)
    tn := dt.ConstTickerNode{5,1337}
    fmt.Println(tn.Exec(5,inpipes))
    fmt.Println(oddsInStream.StreamDef.Exec(6))
}
