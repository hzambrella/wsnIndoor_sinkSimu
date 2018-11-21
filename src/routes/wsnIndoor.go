package routes

import (
	"fmt"
	//"fmt"
	"math"
	"strconv"
	//"math/rand"
	//"time"
	"model/sensor"
	//"github.com/gin-gonic/contrib/sessions"
	//"interceptor/auth"
	"tool/gintool"
	//"tool/captcha"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
)

const (
	startLocAnchorPath="/wsnIndoor/startLocAnchor"
	getHigherAnchorHopPath = "/wsnIndoor/getHigherAnchorHop"
	getNormalAnchorHopPath="/wsnIndoor/getNormalAnchorHop"
)

func init() {
	r := gintool.Default()
	//TODO: 改成post 请求
	r.GET(startLocAnchorPath,startLocAnchor)
	r.GET(getHigherAnchorHopPath, getHigherAnchorHop)
	r.GET(getNormalAnchorHopPath, getNormalAnchorHop)
}

// 模拟测试，让系统开始进行锚节点的标定。逻辑是将build_network表中status改为0
func startLocAnchor(c *gin.Context){
	nid, err := VarsValueInt(c, "nid")
	if err != nil {
		log.Warn(err)
		c.JSON(500, paramWrongFormat("nid"))
		return
	}

	c.JSON(200,gin.H{
		"nid":nid,
	})
}

// 模拟测试:获得高级锚节点之间的跳距矩阵以及坐标,用于训练
func getHigherAnchorHop(c *gin.Context) {
	nid, err := VarsValueInt(c, "nid")
	if err != nil {
		log.Warn(err)
		c.JSON(500, paramWrongFormat("nid"))
		return
	}

	db, err := sensor.NewSensorDB()
	if err != nil {
		log.Warn(err)
		c.JSON(500, dbWrong)
		return
	}

	anchorsHigher, err := db.GetXYAnchor(nid,sensor.AnchorTypeHiger)
	if err != nil {
		if err != sensor.AnchorDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}

	anchorsNormal, err := db.GetXYAnchor(nid, sensor.AnchorTypeNormal)
	if err != nil {
		if err != sensor.AnchorDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}
	
	//higher anchor在前面
	anchors:=append(anchorsHigher,anchorsNormal...)
	
	anchorRadius,err:=db.GetAnchorRadius(nid)
	if err != nil {
		if err != sensor.NetworkDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}

	hop,err:=getHop(anchors,anchorRadius.AnchorRadius)
	if err!=nil{
		log.Warn(err)
		c.JSON(500, sysWrong)
		return
	}

	var higherHop [][]float64
	higherNum:=len(anchorsHigher)
	higherHop = make([][]float64, higherNum)
	for i := 0; i < higherNum; i++ {
		higherHop[i] = make([]float64, higherNum)
	}

	for i:=0;i<higherNum;i++{
		for j:=0;j<higherNum;j++{
			higherHop[i][j]=hop[i][j]
		}
	}

	c.JSON(200, gin.H{
		"anchors":anchorsHigher,
		//"distance":distance,
		"hop":higherHop,
	})
}


// 模拟测试:获得普通锚节点和高级锚节点之间的跳距矩阵,用于获取结果
func getNormalAnchorHop(c *gin.Context) {
		nid, err := VarsValueInt(c, "nid")
	if err != nil {
		log.Warn(err)
		c.JSON(500, paramWrongFormat("nid"))
		return
	}

	db, err := sensor.NewSensorDB()
	if err != nil {
		log.Warn(err)
		c.JSON(500, dbWrong)
		return
	}

	anchorsHigher, err := db.GetXYAnchor(nid,sensor.AnchorTypeHiger)
	if err != nil {
		if err != sensor.AnchorDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}

	anchorsNormal, err := db.GetXYAnchor(nid, sensor.AnchorTypeNormal)
	if err != nil {
		if err != sensor.AnchorDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}
	
	//higher anchor在前面
	anchors:=append(anchorsHigher,anchorsNormal...)
	
	anchorRadius,err:=db.GetAnchorRadius(nid)
	if err != nil {
		if err != sensor.NetworkDataNotFound {
			log.Warn(err)
			c.JSON(500, dbWrong)
			return
		}
	}

	hop,err:=getHop(anchors,anchorRadius.AnchorRadius)
	if err!=nil{
		log.Warn(err)
		c.JSON(500, sysWrong)
		return
	}

	var normalHop [][]float64
	higherNum:=len(anchorsHigher)
	normalNum:=len(anchorsNormal)
	normalHop = make([][]float64, normalNum)
	for i := 0; i < normalNum; i++ {
		normalHop[i] = make([]float64, higherNum)
	}

	fmt.Println(higherNum);
	for i:=0;i<normalNum;i++{
		for j:=0;j<higherNum;j++{
			normalHop[i][j]=hop[higherNum+i][j]
		}
	}

	c.JSON(200, gin.H{
		//"distance":distance,
		"hop":normalHop,
	})
}

func getDistance(anchors []sensor.Anchor) ([][]float64, error) {
	length := len(anchors)
	var disMetrix [][]float64
	disMetrix = make([][]float64, length)
	for i := 0; i < length; i++ {
		disMetrix[i] = make([]float64, length)
	}

	for i := 0; i < length; i++ {
		x1, err := strconv.ParseFloat(anchors[i].X, 64)
		if err != nil {
			return nil, err
		}

		y1, err := strconv.ParseFloat(anchors[i].Y, 64)
		if err != nil {
			return nil, err
		}

		for j := i + 1; j < length; j++ {

			x2, err := strconv.ParseFloat(anchors[j].X, 64)
			if err != nil {
				return nil, err
			}

			y2, err := strconv.ParseFloat(anchors[j].Y, 64)
			if err != nil {
				return nil, err
			}

			distance := math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))

			disMetrix[i][j] = distance
			disMetrix[j][i] = distance
		}
	}

	return disMetrix, nil
}

func getHop(anchors []sensor.Anchor,radius float64)([][]float64,error){
	length := len(anchors)
/*
	var disMetrix [][]float64
	disMetrix = make([][]float64, length)
	for i := 0; i < length; i++ {
		disMetrix[i] = make([]float64, length)
	}
*/
	var hop [][]float64
	hop = make([][]float64, length)
	for i := 0; i < length; i++ {
		hop[i] = make([]float64, length)
	}

	for i:=0;i<length;i++{
		for j:=0;j<length;j++{
			distance,err:=cacuDistanceInAnchor(anchors,i,j)
			if err != nil {
				return nil, err
			}
			//disMetrix[i][j] = distance
			if distance<radius&&distance>0{
				hop[i][j]=1
			}else if i==j{
				hop[i][j]=0
			}else{
				hop[i][j]=math.Inf(1)
			}
		}
	}

	for k:=0;k<length;k++{
		for i:=0;i<length;i++{
			for j:=0;j<length;j++{
				if hop[i][k]+hop[k][j]<hop[i][j]{
					hop[i][j]=hop[i][k]+hop[k][j]
				}
			}
		}
	}
	
	return hop,nil
}

func cacuDistanceInAnchor(anchors []sensor.Anchor,i int,j int )(float64,error){
	x1, err := strconv.ParseFloat(anchors[i].X, 64)
	if err != nil {
		return 0, err
	}

	y1, err := strconv.ParseFloat(anchors[i].Y, 64)
	if err != nil {
		return 0, err
	}

	x2, err := strconv.ParseFloat(anchors[j].X, 64)
	if err != nil {
		return 0, err
	}

	y2, err := strconv.ParseFloat(anchors[j].Y, 64)
	if err != nil {
		return 0, err
	}

	distance := math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
	return distance,nil
}