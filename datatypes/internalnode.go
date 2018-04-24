package datatypes

// interface
type internalNode interface {
    Vote (t Time) *Time;
    Exec (t Time, inpipes InPipes) EvPayload;
    Rinse (inpipes InPipes)
}

// const ticker

func (node ConstTickerNode) Vote (t Time) *Time {
    if t<node.Value {
        return &node.Value
    }
    return nil
}

func (node ConstTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
    if t==node.Value {
        return UnitPayload
    }
    return NothingPayload
}

func (node ConstTickerNode) Rinse (inpipes InPipes) {
}

// src ticker

func (node SrcTickerNode) Vote (t Time) *Time {
    return nil
}

func (node SrcTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
    ev,ok := inpipes.consume(SrcSignal)
    if !ok {
        panic("No input signal!")
    }
    if ev != nil {
        return UnitPayload
    }
    return NothingPayload
}

func (node SrcTickerNode) Rinse (inpipes InPipes) {
    inpipes.rinse(SrcSignal)
}

// TODO: Union ticker, delay ticker
