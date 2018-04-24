package datatypes

/*
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
*/

// TNode

func (node TNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return some(t)
}

func (node TNode) Rinse (inpipes InPipes) {
}

// WNode

func (node WNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return some(w)
}

func (node WNode) Rinse (inpipes InPipes) {
}

// aux funs

func consumeWhile(seen []Event, cmpfun func(Time) bool) ([]Event, EvPayload) {
    if (len(seen) == 0 || !cmpfun(seen[0].time)) {
        return seen, NothingPayload // outside
    }
    i:=1
    for ;i<len(seen);i++ {
        if !cmpfun(seen[i].time) {
            break
        }
    }
    seen = seen[i-1:]
    return seen, some(seen[0])
}

// PrevValNode

func (node PrevValNode) Exec (_ Time, _ interface{}, inpipes InPipes) EvPayload {
    ev,ok := inpipes.consume(TickSignal)
    if !ok {
        panic("No input tick signal!")
    }
    if !ev.payload.isSet {
        // outside
        return ev.payload
    }
    limitT := ev.payload.val.(Time)
    lowerthant := func(seent Time) bool {
        return seent<limitT
    }
    newseen, rett := consumeWhile(node.seen, lowerthant)
    node.seen = newseen
    return rett
}

func (node PrevValNode) Rinse (inpipes InPipes) {
    ev,ok := inpipes.consume(SrcSignal)
    if !ok {
        panic("No input signal!")
    }
    if (ev !=nil) {
        node.seen = append(node.seen, *ev)
    }
}

// PrevEqValNode

func (node PrevEqValNode) Exec (t Time, _ interface{}, inpipes InPipes) EvPayload {
    ev,ok := inpipes.consume(TickSignal)
    if !ok {
        panic("No input tick signal!")
    }
    if !ev.payload.isSet {
        // outside
        return ev.payload
    }
    limitT := ev.payload.val.(Time)
    if (limitT == t) {
        // It might be now
        node.Rinse(inpipes)
    }
    leqthant := func(seent Time) bool {
        return seent<=limitT
    }
    newseen, rett := consumeWhile(node.seen, leqthant)
    node.seen = newseen
    return rett
}

func (node PrevEqValNode) Rinse (inpipes InPipes) {
    ev,ok := inpipes.consume(SrcSignal)
    if !ok {
        panic("No input signal!")
    }
    if (ev !=nil) {
        node.seen = append(node.seen, *ev)
    }
}
