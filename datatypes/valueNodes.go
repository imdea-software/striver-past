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

// Beta testing: generic funs

func genericExec(t Time, w interface{}, inpipes InPipes, rinsefun func(InPipes), tpointernode ValNode, cmpfun func(t0 Time, t1 Time) bool, seen *[]Event, updateSeen func([]Event), extractor func(Event) interface{}) EvPayload {
    tpayload := tpointernode.Exec(t,w,inpipes)
    if !tpayload.IsSet {
        // outside
        return tpayload
    }
    limitT := tpayload.Val.(Time)
    if (limitT == t) {
        // It might be now
        rinsefun(inpipes)
    }
    newseen, rett := consumeWhile(*seen, func(t Time) bool {return cmpfun(t, limitT)})
    updateSeen(newseen)
    if !rett.IsSet {
        // outside
        return rett
    }
    ev := rett.Val.(Event)
    if (!ev.Payload.IsSet) {
        panic("Empty payload in queue??")
    }
    return Some(extractor(ev))
}

func genericRinse (inpipes InPipes, tpointernode ValNode, srcStream StreamName, seen []Event, updateSeen func([]Event)) {
    tpointernode.Rinse(inpipes)
    ev := inpipes.strictConsume(srcStream)
    if ev.Payload.IsSet && (len(seen) == 0 || ev.Time!=seen[len(seen)-1].Time) {
        updateSeen(append(seen, ev))
    }
}

func extractPayload(ev Event) interface{} {
    return ev.Payload
}

func extractTime(ev Event) interface{} {
    return ev.Time
}

// PrevEqValNode

func (node *PrevEqValNode) updateSeen (newseen []Event) {
    node.Seen = newseen
}

func (node *PrevEqValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return genericExec(t, w, inpipes, node.Rinse, node.TPointer, Leq, &node.Seen, node.updateSeen, extractPayload)
}

func (node *PrevEqValNode) Rinse (inpipes InPipes) {
    genericRinse(inpipes, node.TPointer, node.SrcStream, node.Seen, node.updateSeen)
}

// PrevValNode

func (node *PrevValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
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
    if !rett.IsSet {
        // outside
        return rett
    }
    ev := rett.Val.(Event)
    return ev.Payload
}

func (node *PrevValNode) Rinse (inpipes InPipes) {
    node.TPointer.Rinse(inpipes)
    ev := inpipes.strictConsume(node.SrcStream)
    if ev.Payload.IsSet {
        node.Seen = append(node.Seen, ev)
    }
}
