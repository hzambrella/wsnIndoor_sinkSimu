package ginIOC_param

import (
	"fmt"
	"strconv"
	"testing"
	"os"
)

func TestGaoji(t *testing.T) {
	var pkg string = os.Getenv("PJDIR") + "/routes"
	IOCParam(pkg)
}

func TestStrconv(t *testing.T) {
	btStr := "1213"
	if btStr != "" {
		fmt.Println(1)
	}

	bt, _ := strconv.ParseUint(btStr, 10, 32)
	fmt.Println(bt)
	if err := fa(); err != nil {

	}
}

func fa() error {
	return nil
}
