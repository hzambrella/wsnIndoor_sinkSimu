package pgdatastore

import (
	"os"
	"database/sql"
	"tool/inicfg"
	"sync"

	_ "github.com/lib/pq"
)

//前面的杠表示，运行这个包的init函数。
//参考http://www.tuicool.com/articles/jyqq63
//	_ "github.com/go-sql-driver/mysql"


type ConnectDB interface {
	GetDB(dbName string) (*sql.DB,error)
}

//连接数据库对系统而言是很大的开销。缓存连接的目的是减少开销。这里模仿 柳丁的git.ot24.net/go/engine/datastore包,可能不周全
var mu sync.Mutex
var LinkStore DBstore

type DBstore struct {
	Ms map[string]*sql.DB
}

func init(){
	m := make(map[string]*sql.DB, 0)
	LinkStore = DBstore{m}
}

//获取数据库连接
func (d DBstore) GetDB(linkName string) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	cfg:=inicfg.Getcfg()
	if cfg==nil{
		var errCfg error
		cfg,errCfg=inicfg.Newcfg(os.Getenv("ETCDIR"))
		if errCfg!=nil{
			return nil,errCfg
		}
	}

    mapDBCfg,err:=inicfg.Getcfg().GetSection(linkName)
	if err!=nil{
		return nil,err
	}

	db, ok := d.Ms[linkName]
	if !ok || db == nil {
		driver:=mapDBCfg["driver"]
		dsnCfg:=mapDBCfg["dsn"]
		newlink, err := sql.Open(driver, dsnCfg)
		if err != nil {
			return nil, err
		}
		d.Ms[linkName] = newlink
		return newlink, nil
	}
	return db, nil
}

