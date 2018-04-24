package datatypes

// interface
type internalNode interface {
    Exec (t Time, inpipes InPipes) EvPayload;
    Rinse (t Time, inpipes InPipes)
}

type TickerNode interface {
    Vote (t Time) *Time;
}

// tickers
type ConstTickerNode struct {
    ConstT Time
    ConstW Time
}

type SrcTickerNode struct {
}

type DelayTickerNode struct {
    Combiner func(a EvPayload, b EvPayload) EvPayload;
    alarms []Event
}

type UnionTickerNode struct {
    Combiner func(a EvPayload, b EvPayload) EvPayload
}

// values

type TNode struct {
}

type WNode struct {
}

type PrevNode struct {
    seen []Time
}

type PrevEqNode struct {
    seen []Time
}

type PrevValNode struct {
    seen []Event
}

type PrevEqValNode struct {
    seen []Event
}

type FuncNode struct {
    Innerfun func (args ...EvPayload) EvPayload
}

// outs

type CoreNode struct {
}

type InNode struct {
}
