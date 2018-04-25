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
    return Some(t)
}

func (node TNode) Rinse (_ InPipes) {
}

// WNode

func (node WNode) Exec (_ Time, w interface{}, _ InPipes) EvPayload {
    return Some(w)
}

func (node WNode) Rinse (inpipes InPipes) {
}

// aux funs

func consumeWhile(seen []Event, cmpfun func(Time) bool) ([]Event, EvPayload) {
    if (len(seen) == 0 || !cmpfun(seen[0].Time)) {
        return seen, NothingPayload // outside
    }
    i:=1
    for ;i<len(seen);i++ {
        if !cmpfun(seen[i].Time) {
            break
        }
    }
    seen = seen[i-1:]
    return seen, Some(seen[0])
}

// PrevValNode

func (node PrevValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    tpayload := node.TPointer.Exec(t,w,inpipes)
    if !tpayload.IsSet {
        // outside
        return tpayload
    }
    limitT := tpayload.Val.(Time)
    lowerthant := func(seent Time) bool {
        return seent<limitT
    }
    newseen, rett := consumeWhile(node.Seen, lowerthant)
    node.Seen = newseen
    return rett
}

func (node PrevValNode) Rinse (inpipes InPipes) {
    node.TPointer.Rinse(inpipes)
    ev := inpipes.strictConsume(node.SrcStream)
    if ev.Payload.IsSet {
        node.Seen = append(node.Seen, ev)
    }
}

// PrevEqValNode

func (node PrevEqValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    tpayload := node.TPointer.Exec(t,w,inpipes)
    if !tpayload.IsSet {
        // outside
        return tpayload
    }
    limitT := tpayload.Val.(Time)
    if (limitT == t) {
        // It might be now
        ev := inpipes.strictConsume(node.SrcStream)
        if ev.Payload.IsSet {
            node.Seen = append(node.Seen, ev)
        }
    }
    leqthant := func(seent Time) bool {
        return seent<=limitT
    }
    newseen, rett := consumeWhile(node.Seen, leqthant)
    node.Seen = newseen
    if !rett.IsSet {
        // outside
        return rett
    }
    ev := rett.Val.(Event)
    return ev.Payload
}

func (node PrevEqValNode) Rinse (inpipes InPipes) {
    node.TPointer.Rinse(inpipes)
    ev := inpipes.strictConsume(node.SrcStream)
    if ev.Payload.IsSet && (len(node.Seen) == 0 || ev.Time!=node.Seen[len(node.Seen)-1].Time) {
        node.Seen = append(node.Seen, ev)
    }
}
