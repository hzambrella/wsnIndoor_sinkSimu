package datastore

import(
	"github.com/Unknwon/goconfig"
	"os"
	"tool/inicfg"
	"testing"
)
var cfg *goconfig.ConfigFile

func init() {
	var err error
    cfg,err=inicfg.Newcfg(os.Getenv("ETCDIR"));
	if err!=nil{
		panic(err)
	}
}

func TestEngine(t *testing.T) {
	db,err:=LinkStore.GetDB("master")
	if err!=nil{
		t.Fatal(err)
	}
	if err:=db.Ping();err!=nil{
		t.Fatal(err)
	}
}