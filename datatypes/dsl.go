package datatypes

type Time int

type StreamName string

type Event struct {
    Time Time;
    Payload EvPayload
}

type EvPayload struct {
    IsSet bool;
    Val interface{}
}

func Some(val interface{}) EvPayload {
    return EvPayload{true, val}
}

var NothingPayload EvPayload = EvPayload{false, nil}

type EpsVal struct {
    eps Time
    val interface{}
}

type OutStream struct {
    Name StreamName
    TicksDef TickerNode
    ValDef ValNode
}

type InStream struct {
    Name StreamName
    StreamDef InStreamDef
}

type InStreamDef interface {
    PeekNextTime() Time // infinite input streams
    Exec(t Time) EvPayload
}
