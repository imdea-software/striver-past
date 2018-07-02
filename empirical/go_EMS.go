package empirical

import (
    "fmt"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "math/rand"
)

func ArrivalStock(products int) (inStreams []dt.InStream, outStreams []dt.OutStream, killcallback func()) {
    // Input streams
    saleName := dt.StreamName("sale")
    saleChan := make(chan dt.Event)
    arrivalName := dt.StreamName("arrival")
    arrivalChan := make(chan dt.Event)
	sale := dt.InStream{saleName, &dt.InFromChannel{saleChan, nil, 0, false}}
	arrival := dt.InStream{arrivalName, &dt.InFromChannel{arrivalChan, nil, 0, false}}
    inStreams = []dt.InStream{sale, arrival}

    // Output stream:
    stockName := dt.StreamName("stock")
    // ticks
    saleSrcTick := dt.SrcTickerNode{saleName}
    arrivalSrcTick := dt.SrcTickerNode{arrivalName}
    stockTicks := dt.UnionTickerNode{saleSrcTick, arrivalSrcTick, dt.FstPayload}
    // val
    tpointer := dt.TNode{}
    stockPrevVal := dt.PrevValNode{tpointer,stockName, []dt.Event{}}
    arrivalPrevEq := dt.PrevEqNode{tpointer, arrivalName, []dt.Event{}}
    arrivalPrevEqVal := dt.PrevEqValNode{tpointer, arrivalName, []dt.Event{}}
    salePrevEq := dt.PrevEqNode{tpointer, saleName, []dt.Event{}}
    salePrevEqVal := dt.PrevEqValNode{tpointer, saleName, []dt.Event{}}
    stockFun := func(args ...dt.EvPayload) dt.EvPayload{
            t := args[0]
            lastStock := args[1]
            preveqArrival := args[2]
            preveqvalArrival := args[3]
            preveqSale := args[4]
            preveqvalSale := args[5]

            stock := 0
            if lastStock.IsSet {
                stock = lastStock.Val.(dt.EvPayload).Val.(int)
            }
            if preveqArrival.IsSet && preveqArrival.Val == t.Val {
                // Arrival is ticking
                stock += preveqvalArrival.Val.(dt.EvPayload).Val.(int)
            }
            if preveqSale.IsSet && preveqSale.Val == t.Val {
                // Arrival is ticking
                stock -= preveqvalSale.Val.(dt.EvPayload).Val.(int)
            }
            return dt.Some(stock)
        }
    stockVal := dt.FuncNode{[]dt.ValNode{tpointer, &stockPrevVal, &arrivalPrevEq, &arrivalPrevEqVal, &salePrevEq, &salePrevEqVal}, stockFun}

    stock := dt.OutStream{stockName, stockTicks, stockVal}

    outStreams = []dt.OutStream{stock}

    evCountArrival := 0
    evCountSale := 0
    killcallback = func() { fmt.Println("Processed events:", evCountArrival + evCountSale) }

    // Feed data
    go func() {
        nextArrival := 500 + rand.Int63n(20) + 1
        for {
            nextArrival = nextArrival + rand.Int63n(20) + 1
            arrivalChan <- dt.Event{dt.Time(nextArrival), dt.Some(10)}
            evCountArrival++
        }
        close(arrivalChan)
    }()
    go func() {
        nextSale := 500 + rand.Int63n(20) + 1
        for {
            nextSale = nextSale + rand.Int63n(20) + 1
            saleChan <- dt.Event{dt.Time(nextSale), dt.Some(15)}
            evCountSale++
        }
        close(saleChan)
    }()
    return
}

