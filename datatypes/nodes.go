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
