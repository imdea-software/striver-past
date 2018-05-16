package datatypes

// interface
type TickerNode interface {
    Vote (t Time) MaybeTime;
    Exec (t Time, inpipes InPipes) EvPayload;
    Rinse (inpipes InPipes)
}

type ValNode interface {
    Exec (t Time, w interface{}, inpipes InPipes) EvPayload;
    Rinse (inpipes InPipes)
}

// tickers
type ConstTickerNode struct {
    ConstT Time
    ConstW interface{}
}

type SrcTickerNode struct {
    SrcStream StreamName
}

type DelayTickerNode struct {
    SrcStream StreamName
    Combiner func(a EvPayload, b EvPayload) EvPayload;
    Alarms []Event
}

type UnionTickerNode struct {
    LeftTicker TickerNode
    RightTicker TickerNode
    Combiner func(a EvPayload, b EvPayload) EvPayload
}

// values

type TNode struct {
}

type WNode struct {
}

type PrevNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevEqNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevValNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type PrevEqValNode struct {
    TPointer ValNode
    SrcStream StreamName
    Seen []Event
}

type FuncNode struct {
    ArgNodes []ValNode
    Innerfun func (args ...EvPayload) EvPayload
}

type InFromChannel struct {
    InChannel chan Event
    NextEvent *Event
}

func (ticker *InFromChannel) PeekNextTime () MaybeTime {
    if ticker.NextEvent == nil {
        nextEv := <-ticker.InChannel
        ticker.NextEvent = &nextEv
    }
    return SomeTime(ticker.NextEvent.Time)
}

func (ticker *InFromChannel) Exec (t Time) EvPayload {
    if t == ticker.NextEvent.Time {
        ret := Some(ticker.NextEvent.Payload)
        ticker.NextEvent = nil
        return ret
    }
    return NothingPayload
}
