package sensor

import (
	//"tool/errors"
	//"database/sql"
	//"encoding/json"
	"fmt"
	"testing"
	//"time"
)

var dbu SensorDB

func init() {
	var err error

	dbu, err = NewSensorDB()
	if err != nil {
		panic(err)
	}
}

func TestGetXY(t *testing.T){
	anchors,err:=dbu.GetXYAnchor(1,AnchorTypeHiger)
	if err != nil {
		if err == AnchorDataNotFound {
			fmt.Println(AnchorDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(anchors)
}

func TestGetAnchorRadius(t *testing.T){
	anchorRadius,err:=dbu.GetAnchorRadius(1);
	if err != nil {
		if err == AnchorDataNotFound {
			fmt.Println(AnchorDataNotFound.Error())
			return
		}
		t.Fatal(err)
	}
	fmt.Println(anchorRadius)
}