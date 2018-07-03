package main

import (
    "fmt"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "os/signal"
    "os"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/empirical"
)

func main() {

    var lastEvent dt.FlowingEvent
    // inStreams, outStreams, killcb := empirical.ArrivalStock(100)
    inStreams, outStreams, killcb := empirical.LastK(100)

    //inStreams, outStreams := shiftExample()
    // inStreams, outStreams := changePointsExample()
    //inStreams, outStreams := clockExample()
    kchan := make (chan bool)
    outchan := make (chan dt.FlowingEvent)
    go func(){
        for ev := range outchan {
            lastEvent = ev
            // Ignore incoming events
        }
    }()

    // Catch interruption
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func(){
        for _ = range c {
            close(kchan)
        }
    }()

    controlplane.Start(inStreams, outStreams, outchan, kchan)
    fmt.Println("End of execution")
    killcb()
    fmt.Println(lastEvent)
}
