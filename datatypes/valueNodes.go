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

func (node TNode) Exec (t Time, _ interface{}, _ InPipes) EvPayload {
    return some(t)
}

func (node TNode) Rinse (_ InPipes) {
}

// WNode

func (node WNode) Exec (_ Time, w interface{}, _ InPipes) EvPayload {
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

func (node PrevValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    tpayload := node.TPointer.Exec(t,w,inpipes)
    if !tpayload.isSet {
        // outside
        return tpayload
    }
    limitT := tpayload.val.(Time)
    lowerthant := func(seent Time) bool {
        return seent<limitT
    }
    newseen, rett := consumeWhile(node.seen, lowerthant)
    node.seen = newseen
    return rett
}

func (node PrevValNode) Rinse (inpipes InPipes) {
    node.TPointer.Rinse(inpipes)
    ev := inpipes.strictConsume(node.SrcStream)
    if ev.payload.isSet {
        node.seen = append(node.seen, ev)
    }
}

// PrevEqValNode

func (node PrevEqValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    tpayload := node.TPointer.Exec(t,w,inpipes)
    if !tpayload.isSet {
        // outside
        return tpayload
    }
    limitT := tpayload.val.(Time)
    if (limitT == t) {
        // It might be now
        ev := inpipes.strictConsume(node.SrcStream)
        if ev.payload.isSet {
            node.seen = append(node.seen, ev)
        }
    }
    leqthant := func(seent Time) bool {
        return seent<=limitT
    }
    newseen, rett := consumeWhile(node.seen, leqthant)
    node.seen = newseen
    return rett
}

func (node PrevEqValNode) Rinse (inpipes InPipes) {
    node.TPointer.Rinse(inpipes)
    ev := inpipes.strictConsume(node.SrcStream)
    if ev.payload.isSet && (len(node.seen) == 0 || ev.time!=node.seen[len(node.seen)-1].time) {
        node.seen = append(node.seen, ev)
    }
}
