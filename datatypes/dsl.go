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

func some(val interface{}) EvPayload {
    return EvPayload{true, val}
}

var NothingPayload EvPayload = EvPayload{false, nil}

type EpsVal struct {
    eps Time
    val interface{}
}