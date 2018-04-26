package datatypes

func Min(t0 *Time, t1 *Time) *Time {
    if t0 == nil {
        return t1
    }
    if t1 == nil {
        return t0
    }
    if *t0 < *t1 {
        return t0
    }
    return t1
}

func Leq(t0 Time, t1 Time) bool {
    return t0 <= t1
}

func Lt(t0 Time, t1 Time) bool {
    return t0 < t1
}

func FstPayload(t0 EvPayload, _ EvPayload) EvPayload {
    return t0
}
