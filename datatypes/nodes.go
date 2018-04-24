package datatypes

// tickers
type ConstTickerNode struct {
    Value Time
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
