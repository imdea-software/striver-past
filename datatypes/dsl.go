package datatypes

type Time int

type StreamName string

type PipeId int

const (
   SrcSignal    PipeId = 0
)

type Event struct {
    time Time;
    payload EvPayload
}

type EvPayload struct {
    isSet bool;
    val interface{}
}

var UnitPayload EvPayload = EvPayload{true, nil}
var NothingPayload EvPayload = EvPayload{false, nil}
