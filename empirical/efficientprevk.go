package empirical

import (
    "fmt"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "math/rand"
)

func EffLastK(k int) (inStreams []dt.InStream, outStreams []dt.OutStream, killcallback func()) {

    inStreams = []dt.InStream{}
    outStreams = []dt.OutStream{}
    evCount := 0
    // Input streams
    saleName := dt.StreamName("sale")
    saleChan := make(chan dt.Event)
    sale := dt.InStream{saleName, &dt.InFromChannel{saleChan, nil, 0, false}}
    inStreams = append(inStreams, sale)

    // Output stream:
    denomName := dt.StreamName("denom")
    sumKName := dt.StreamName("sumk")
    avgKName := dt.StreamName("avgk")
    // ticks
    denomTicks := dt.SrcTickerNode{saleName}
    sumKTicks := dt.SrcTickerNode{saleName}
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

    sumFun := func(args ...dt.EvPayload) dt.EvPayload {
        lastSum := 0
        rmSale := 0
        if args[0].IsSet {
            lastSum = args[0].Val.(dt.EvPayload).Val.(int)
        }
        sale := args[1].Val.(dt.EvPayload).Val.(int)
        if args[2].IsSet {
            rmSale = args[2].Val.(dt.EvPayload).Val.(int)
        }
        return dt.Some(lastSum + sale - rmSale)
    }
    sumKPrevVal := dt.PrevValNode{tpointer,sumKName, []dt.Event{}}
    salePrevEqVal := dt.PrevEqValNode{tpointer,saleName, []dt.Event{}}
    var lastPointer dt.ValNode = tpointer
    for i:=0 ; i<k ; i++ {
        lastPointer = &dt.PrevNode{lastPointer, saleName, []dt.Event{}}
    }
    salePrevEqValK := dt.PrevEqValNode{lastPointer, saleName, []dt.Event{}}
    sumKVal := dt.FuncNode{[]dt.ValNode{&sumKPrevVal, &salePrevEqVal, &salePrevEqValK}, sumFun}
    sumK := dt.OutStream{sumKName, sumKTicks, sumKVal}
    outStreams = append(outStreams, sumK)

    avgFun := func(args ...dt.EvPayload) dt.EvPayload {
        numer := args[0].Val.(dt.EvPayload).Val.(int)
        denom := args[1].Val.(dt.EvPayload).Val.(int)
        return dt.Some(numer/denom)
    }
    sumPrevEqVal := dt.PrevEqValNode{tpointer, sumKName, []dt.Event{}}
    denomPrevEqVal := dt.PrevEqValNode{tpointer,denomName, []dt.Event{}}
    avgKVal := dt.FuncNode{[]dt.ValNode{&sumPrevEqVal, &denomPrevEqVal}, avgFun}
    avgK := dt.OutStream{avgKName, avgKTicks, avgKVal}
    outStreams = append(outStreams, avgK)

    killcallback = func() { fmt.Println("Processed events:", evCount) }

    // Feed data
    go func() {
        nextev := 500 + rand.Int63n(20) + 1
        for {
            nextev = nextev + rand.Int63n(20) + 1
            saleChan <- dt.Event{dt.Time(nextev), dt.Some(rand.Intn(20))}
            evCount++
            //fmt.Println("sending ev ",nextev)
        }
        close(saleChan)
    }()
    return
}
