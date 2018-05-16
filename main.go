package main

import (
    "fmt"
    "gitlab.software.imdea.org/felipe.gorostiaga/striver-go/controlplane"
)

func main() {

    inStreams, outStreams := shiftExample()
    // inStreams, outStreams := changePointsExample()
    //inStreams, outStreams := clockExample()
    controlplane.Start(inStreams, outStreams)

    fmt.Println("End of execution")
}
