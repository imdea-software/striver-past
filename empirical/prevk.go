package empirical

import (
    //"fmt"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "math/rand"
)

func LastK(k int) (inStreams []dt.InStream, outStreams []dt.OutStream, killcallback func()) {

    inStreams = []dt.InStream{}
    outStreams = []dt.OutStream{}
    // Input streams
    saleName := dt.StreamName("sale")
    saleChan := make(chan dt.Event)
    sale := dt.InStream{saleName, &dt.InFromChannel{saleChan, nil, 0, false}}
    inStreams = append(inStreams, sale)

    // Output stream:
    denomName := dt.StreamName("denom")
    avgKName := dt.StreamName("avgk")
    // ticks
    denomTicks := dt.SrcTickerNode{saleName}
    avgKTicks := dt.SrcTickerNode{saleName}
    // val
    denomFun := func(args ...dt.EvPayload) dt.EvPayload {
        lastDenom := args[0]
        denom := 0
        if lastDenom.IsSet {
            denom = lastDenom.Val.(dt.EvPayload).Val.(int)
        }
        if denom<k {
            denom++
        }
        return dt.Some(denom)
    }
    tpointer := dt.TNode{}
    denomPrevVal := dt.PrevValNode{tpointer,denomName, []dt.Event{}}
    denomVal := dt.FuncNode{[]dt.ValNode{&denomPrevVal}, denomFun}
    denom := dt.OutStream{denomName, denomTicks, denomVal}
    outStreams = append(outStreams, denom)

    avgFun := func(args ...dt.EvPayload) dt.EvPayload {
        denom := args[0].Val.(dt.EvPayload).Val.(int)
        sum := 0
        for i:=1;i<len(args);i++ {
            if args[i].IsSet {
                sum += args[i].Val.(dt.EvPayload).Val.(int)
            //} else {
            //    break
            }
        }
        return dt.Some(sum/denom)
    }

    denomPrevEqVal := dt.PrevEqValNode{tpointer,denomName, []dt.Event{}}
    args := []dt.ValNode{&denomPrevEqVal}

    for i:=0 ; i<k ; i++ {
        var lastPointer dt.ValNode = tpointer
        for j:=0; j<i; j++ {
            lastPointer = &dt.PrevNode{lastPointer, saleName, []dt.Event{}}
        }
        salePrevEqVal := dt.PrevEqValNode{lastPointer, saleName, []dt.Event{}}
        args = append(args, &salePrevEqVal)
    }
    avgKVal := dt.FuncNode{args, avgFun}
    avgK := dt.OutStream{avgKName, avgKTicks, avgKVal}
    outStreams = append(outStreams, avgK)
    killcallback = func() {}

    // Feed data
    go func() {
        nextev := 500 + rand.Int63n(20) + 1
        for {
            nextev = nextev + rand.Int63n(20) + 1
            saleChan <- dt.Event{dt.Time(nextev), dt.Some(rand.Intn(20))}
            //fmt.Println("sending ev ",nextev)
        }
        close(saleChan)
    }()
    return
}
