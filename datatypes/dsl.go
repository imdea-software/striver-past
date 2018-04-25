package datatypes

type Time int

type StreamName string

type Event struct {
    time Time;
    payload EvPayload
}

type EvPayload struct {
    isSet bool;
    val interface{}
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
    name StreamName
    ticksDef TickerNode
    valDef ValNode
}

type InStream struct {
    Name StreamName
    StreamDef InStreamDef
}

type InStreamDef interface {
    PeekNextTime() Time // infinite input streams
    Exec(t Time) EvPayload
}
