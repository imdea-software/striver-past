package main

import (
    "fmt"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "os/signal"
    "os"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/empirical"
    "strconv"
    "time"
)

func main() {

    var lastEvent dt.FlowingEvent
    //events := []dt.FlowingEvent{}
    // inStreams, outStreams, killcb := empirical.ArrivalStock(100)
    arg2,err := strconv.Atoi(os.Args[2])
    if err != nil {
        panic(err)
    }
    maxevs,err := strconv.Atoi(os.Args[3])
    if err != nil {
        panic(err)
    }
    inStreams, outStreams, killcb := empirical.ArrivalStock(arg2,maxevs)
    if os.Args[1]=="AVGK" {
        fmt.Fprintf(os.Stderr, "Running AVGK with K=%d",arg2)
        inStreams, outStreams, killcb = empirical.EffLastK(arg2,maxevs)
    } else {
        fmt.Fprintf(os.Stderr, "Running STOCK with P=%d",arg2)
    }
    fmt.Fprintf(os.Stderr, ", maxevs=%d\n",maxevs)

    //inStreams, outStreams := shiftExample()
    // inStreams, outStreams := changePointsExample()
    //inStreams, outStreams := clockExample()
    kchan := make (chan bool)
    outchan := make (chan dt.FlowingEvent)
    go func(){
        for ev := range outchan {
            lastEvent = ev
            //events = append(events, ev)
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

    start := time.Now()
    controlplane.Start(inStreams, outStreams, outchan, kchan)
    elapsed := time.Since(start)
    //fmt.Fprintf(os.Stderr, "Took %f seconds\n", elapsed.Seconds())
    fmt.Fprintf(os.Stderr, "Event ratio is %f events per second\n", float64(maxevs)/(elapsed.Seconds()))
    fmt.Println("End of execution")
    killcb()
    fmt.Println(lastEvent)
    //fmt.Println(events[0])
}
