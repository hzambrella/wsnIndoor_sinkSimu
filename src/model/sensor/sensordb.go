package sensor

import (
	"database/sql"
	//"encoding/json"

	"fmt"
	//"strconv"
	//"strings"
	//"time"

	"tool/errors"
	"tool/pgdatastore"
	//"tool/money"

	//"github.com/satori/go.uuid"
	//"github.com/qiniu/log"
	//"tool/bcrypt"
)

var(
	AnchorDataNotFound=errors.New("anchor data not found")
	NetworkDataNotFound=errors.New("network data not found")
)

type SensorDB interface {
	Close()
	GetXYAnchor(bid int,nid int,anchorType string) ([]Anchor, error)
	GetAnchorRadius(nid int,floor int)(*AnchorRadius,error)
	UpdateNetworkStatus(bid,floor int,status int)error
}

type sensorDB struct {
	*sql.DB
}

func NewSensorDB() (SensorDB, error) {
	db, err := pgdatastore.LinkStore.GetDB("master")
	if err != nil {
		return nil, errors.As(err)
	}
	udb := &sensorDB{db}
	return udb, nil
}

func (db *sensorDB) Close() {
	db.Close()
}

func (db *sensorDB)GetAnchorRadius(nid int,floor int)(*AnchorRadius,error){
	anchorRadius:=AnchorRadius{};
	err:=db.QueryRow(getAnchorRadiusSql,nid,floor).Scan(
		&anchorRadius.Nid,
		&anchorRadius.AnchorRadius,
		&anchorRadius.Floor,
	)

	if err!=nil{
		if err == sql.ErrNoRows {
			return nil, errors.As(NetworkDataNotFound)
		} else {
			return nil, errors.As(err)
		}
	}
	return &anchorRadius,nil
}

func (db *sensorDB) GetXYAnchor(bid int,nid int,anchorType string) ([]Anchor, error) {
	anchors := []Anchor{}

	getXYOfAnchorSql:=fmt.Sprintf(getXYOfAnchorSqlTemp,bid,nid)
	rows, err := db.Query(getXYOfAnchorSql,anchorType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.As(AnchorDataNotFound)
		} else {
			return nil, errors.As(err)
		}
	}

	for rows.Next(){
		anchor:=Anchor{}
		if err:=rows.Scan(
			&anchor.Gid,
			&anchor.Type,
			&anchor.X,
			&anchor.Y,
		);err!=nil{
			if err == sql.ErrNoRows {
				return nil, errors.As(AnchorDataNotFound)
			} else {
				return nil, errors.As(err)
			}
		}
		anchors=append(anchors,anchor)
	}
	return anchors,nil
}

func (db *sensorDB) UpdateNetworkStatus(bid,floor int,status int)error{
	_,err:=db.Exec(updateNetworkStatusSql,status,bid,floor)
	if err!=nil{
		return errors.As(err)
	}
	return err
}