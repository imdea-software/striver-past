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
        return Some(node.ConstW)
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
    return ev.Payload
}

func (node SrcTickerNode) Rinse (inpipes InPipes) {
}

// delay ticker

func (node *DelayTickerNode) Vote (t Time) *Time {
    if len(node.Alarms)==0 {
        return nil
    }
    return &node.Alarms[0].Time
}

func insertInPlace(alarms []Event, newev Event, combiner func(a EvPayload, b EvPayload) EvPayload) []Event {
    i:=0
    for ;i < len(alarms); i++ {
        if (alarms[i].Time == newev.Time) {
            // Use combiner
            alarms[i].Payload = combiner(alarms[i].Payload, newev.Payload)
            return alarms
        }
        if (alarms[i].Time > newev.Time) {
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

func (node *DelayTickerNode) Exec (t Time, inpipes InPipes) EvPayload {
    if len(node.Alarms)>0 && t==node.Alarms[0].Time {
        ret := Some(node.Alarms[0].Payload)
        node.Alarms = node.Alarms[1:]
        return ret
    }
    return NothingPayload
}

func (node *DelayTickerNode) Rinse (inpipes InPipes) {
    ev := inpipes.strictConsume(node.SrcStream)
    if (ev.Payload.IsSet) {
        payload := ev.Payload.Val.(EpsVal)
        newev := Event{ev.Time+payload.Eps, Some(payload.Val)}
        node.Alarms = insertInPlace(node.Alarms, newev, node.Combiner)
    }
}

// TODO: Union ticker

