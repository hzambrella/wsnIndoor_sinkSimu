package checkinput

import (
	"strconv"
	"fmt"
	"testing"
)

func TestCheckMobile(t *testing.T) {
	fmt.Println(IsPhone("15671882412"))
	fmt.Println(IsPhone("18986884521"))
	fmt.Println(IsPhone("12312321312"))
	fmt.Println(IsPhone("123123213121"))
	fmt.Println(IsPhone("1231232131"))
	fmt.Println(IsPhone("154123121"))
	//TODO：测试其它的函数
}

func TestCheckNumber(t *testing.T) {
	fmt.Println(CheckNumber("15671882412"))
	fmt.Println(CheckNumber("-1"))
	fmt.Println(CheckNumber("-1231231.112321"))
	fmt.Println(CheckNumber("0"))
	fmt.Println(CheckNumber("0000"))
	fmt.Println(CheckNumber("!0.12313"))
	fmt.Println(CheckNumber("aaaaa"))

	s,err:=strconv.Atoi("0000")
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(s)

	
	f,err:=strconv.ParseFloat("0.12313",64)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(f)
	
	
}