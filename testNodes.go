package main

import (
    dt "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/datatypes"
    "math/rand"
    "time"
)

func clockExample() ([]dt.InStream,[]dt.OutStream) {
    constTicker := dt.ConstTickerNode{5, nil}
    delayTicker := dt.DelayTickerNode{"clock", dt.FstPayload, []dt.Event{}}
    constVal := dt.FuncNode{[]dt.ValNode{}, func(...dt.EvPayload) dt.EvPayload{return dt.Some(dt.EpsVal{2, nil})}}
    clock := dt.OutStream{"clock", dt.UnionTickerNode{constTicker, &delayTicker, dt.FstPayload}, constVal}
    return []dt.InStream{}, []dt.OutStream{clock}
}

type randomIntIn struct {
    promisedT dt.Time
    randomOffset dt.Time
    randomMin int
    randomMax int
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func (ticker *randomIntIn) PeekNextTime () dt.MaybeTime {
    return dt.SomeTime(ticker.promisedT)
}

func (ticker *randomIntIn) Exec (t dt.Time) dt.EvPayload {
    if (t == ticker.promisedT) {
        ticker.promisedT = t + dt.Time(random(1, int(ticker.randomOffset)))
        return dt.Some(random(ticker.randomMin, ticker.randomMax))
    }
    return dt.NothingPayload
}

func changePointsExample() ([]dt.InStream,[]dt.OutStream) {
    randindef := randomIntIn{0, 5, 0, 3}
    randin := dt.InStream{"randin", &randindef}

    tpointer := dt.TNode{}
    inprev := dt.PrevValNode{tpointer, "randin", []dt.Event{}}
    inpreveq := dt.PrevEqValNode{tpointer, "randin", []dt.Event{}}

    filtercp := func(args ...dt.EvPayload) dt.EvPayload{
            prevt := args[0]
            preveqt := args[1]
            if prevt.IsSet && prevt.Val == preveqt.Val {return dt.NothingPayload}
            return preveqt
        }
    cpval := dt.FuncNode{[]dt.ValNode{&inprev, &inpreveq}, filtercp}
    changingpoints := dt.OutStream{"changingpoints", dt.SrcTickerNode{"randin"}, cpval}
    return []dt.InStream{randin}, []dt.OutStream{changingpoints}
}

func shiftExample() ([]dt.InStream,[]dt.OutStream) {
    randindef := randomIntIn{0, 5, 0, 3}
    randin := dt.InStream{"randin", &randindef}

    tpointer := dt.TNode{}
    inpreveq := dt.PrevEqValNode{tpointer, "randin", []dt.Event{}}

    constVal := dt.FuncNode{[]dt.ValNode{&inpreveq}, func(args ...dt.EvPayload) dt.EvPayload{return dt.Some(dt.EpsVal{2, args[0].Val})}}
    unitrandin := dt.OutStream{"unitrandin", dt.SrcTickerNode{"randin"}, constVal}
    delayTicker := dt.DelayTickerNode{"unitrandin", dt.FstPayload, []dt.Event{}}
    w := dt.WNode{}
    shifted := dt.OutStream{"shifted", &delayTicker, w}
    return []dt.InStream{randin}, []dt.OutStream{unitrandin, shifted}
}

