package datatypes

// const ticker

func (node ConstTickerNode) Vote (t Time) *Time {
    if t<node.ConstT {
        return &node.ConstT
    }
    return nil
}

func (node ConstTickerNode) Exec (t Time, _ InPipes) EvPayload {
    if t==node.ConstT {
        return some(node.ConstW)
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
    ev := inpipes.strictConsume(node.SrcStream)
    return some(ev.payload)
}

func (node SrcTickerNode) Rinse (inpipes InPipes) {
}

// delay ticker

func (node DelayTickerNode) Vote (t Time) *Time {
    if len(node.alarms)==0 {
        return nil
    }
    return &node.alarms[0].time
}

func insertInPlace(alarms []Event, newev Event, combiner func(a EvPayload, b EvPayload) EvPayload) []Event {
    i:=0
    for ;i < len(alarms); i++ {
        if (alarms[i].time == newev.time) {
            // Use combiner
            alarms[i].payload = combiner(alarms[i].payload, newev.payload)
            return alarms
        }
        if (alarms[i].time > newev.time) {
            break
        }
    }
    alarms = append(alarms, newev)
    if (i != len(alarms)) {
        copy(alarms[i+1:], alarms[i:])
        alarms[i] = newev
    }
    return alarms
}

func (node DelayTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
    if t==node.alarms[0].time {
        node.alarms = node.alarms[1:]
        return some(node.alarms[0].payload)
    }
    return NothingPayload
}

func (node DelayTickerNode) Rinse (inpipes InPipes) {
    ev := inpipes.strictConsume(node.SrcStream)
    payload := ev.payload.val.(EpsVal)
    newev := Event{ev.time+payload.eps, some(payload.val)}
    node.alarms = insertInPlace(node.alarms, newev, node.Combiner)
}

// TODO: Union ticker

