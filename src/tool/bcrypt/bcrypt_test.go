package bcrypt

import (
	"fmt"
	"testing"
)

func TestJiami(t *testing.T){
	o:="a123456A"
	e,salt:=BcryptHash(o)
	
	b:=BcryptMatch(o, e)
	fmt.Println("e:",e)
	fmt.Println("salt:",salt)
	fmt.Println("b:",b)
}

func TestJiaoyan(t *testing.T){
	o:="a123456A"
	db:="$2a$10$mELSUhpS3xejw2GSKMoRCuI6E4dHbpdO0uBHkkGP6Zu1LV5QMpSrq"
	b:=BcryptMatch(o, db)
	fmt.Println(b)
}