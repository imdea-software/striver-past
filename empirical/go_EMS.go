package empirical

import (
    "fmt"
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "math/rand"
    "strconv"
)

func ArrivalStock(products int) (inStreams []dt.InStream, outStreams []dt.OutStream, killcallback func()) {

    inStreams = []dt.InStream{}
    outStreams = []dt.OutStream{}
    chans := make([](chan dt.Event), products*2)
    evCounts := make([]int, products*2)
    for i:=0 ; i<products ; i++ {
        // Input streams
        saleName := dt.StreamName("sale_"+strconv.Itoa(i))
        saleChan := make(chan dt.Event)
        chans[i*2] = saleChan
        sale := dt.InStream{saleName, &dt.InFromChannel{saleChan, nil, 0, false}}
        arrivalName := dt.StreamName("arrival_"+strconv.Itoa(i))
        arrivalChan := make(chan dt.Event)
        chans[i*2+1] = arrivalChan
        arrival := dt.InStream{arrivalName, &dt.InFromChannel{arrivalChan, nil, 0, false}}
        inStreams = append(inStreams, sale)
        inStreams = append(inStreams, arrival)

        // Output stream:
        stockName := dt.StreamName("stock_"+strconv.Itoa(i))
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
        stockFun := func(args ...dt.EvPayload) dt.EvPayload {
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
        outStreams = append(outStreams, stock)
    }

    killcallback = func() { fmt.Println("Processed events:", sum(evCounts)) }

    // Feed data
    for i,c := range chans {
        go func(i int, c chan dt.Event) {
            nextev := 500 + rand.Int63n(20) + 1
            for {
                nextev = nextev + rand.Int63n(20) + 1
                c <- dt.Event{dt.Time(nextev), dt.Some(10)}
                evCounts[i]=evCounts[i]+1
                //fmt.Println("sending ev ",evCounts[i],nextev)
            }
            close(c)
        }(i,c)
    }
    return
}

func sum(ints []int) int {
    ret := 0
    for _,i := range ints {
        ret += i
    }
    return ret
}

