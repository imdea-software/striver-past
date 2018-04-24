package datatypes

// const ticker

func (node ConstTickerNode) Vote (t Time) *Time {
    if t<node.ConstT {
        return &node.ConstT
    }
    return nil
}

func (node ConstTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
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
    ev,ok := inpipes.consume(SrcSignal)
    if !ok {
        panic("No input signal!")
    }
    if ev != nil {
        return some(ev.payload)
    }
    return NothingPayload
}

func (node SrcTickerNode) Rinse (inpipes InPipes) {
    inpipes.rinse(SrcSignal)
}

// delay ticker

func (node DelayTickerNode) Vote (t Time) *Time {
    if len(node.alarms)==0 {
        return nil
    }
    return &node.alarms[0].time
}

func insertInPlace(alarms []Event, newev Event, combiner func(a EvPayload, b EvPayload) EvPayload) {
    i:=0
    for ;i < len(alarms); i++ {
        if (alarms[i].time == newev.time) {
            // Use combiner
            alarms[i].payload = combiner(alarms[i].payload, newev.payload)
            return
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
}

func (node DelayTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
    if t==node.alarms[0].time {
        node.alarms = node.alarms[1:]
        return some(node.alarms[0].payload)
    }
    return NothingPayload
}

func (node DelayTickerNode) Rinse (inpipes InPipes) {
    ev,exists := inpipes.consume(SrcSignal)
    if exists {
        payload := ev.payload.val.(EpsVal)
        newev := Event{ev.time+payload.eps, some(payload.val)}
        insertInPlace(node.alarms, newev, node.Combiner)
    }
}

// TODO: Union ticker

