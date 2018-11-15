package routes

import (
	"math"
	"fmt"
	"model/sensor"
	"testing"
)

func TestFunc(t *testing.T) {
	anchors := make([]sensor.Anchor, 3)
	anchors[0] = sensor.Anchor{Gid: 1, X: "0", Y: "0", Type: sensor.AnchorTypeHiger}
	anchors[1] = sensor.Anchor{Gid: 2, X: "2", Y: "0", Type: sensor.AnchorTypeHiger}
	anchors[2] = sensor.Anchor{Gid: 3, X: "1", Y: "1", Type: sensor.AnchorTypeHiger}

	dm, err := getDistance(anchors)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dm)

	fmt.Println(math.Inf(1)>10000)


	hop,err:=getHop(anchors,3,1.6)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hop)
}

