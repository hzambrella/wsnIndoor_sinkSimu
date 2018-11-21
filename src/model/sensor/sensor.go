package sensor

//anchor的类型,数据表中的layer字段
const (
	AnchorTypeHiger  = "anchor_node_higher"
	AnchorTypeNormal = "anchor_node"
)

//sensor
type Anchor struct {
	Gid  int    //gid
	X    string //x
	Y    string //y
	Type string // 类型
}

type AnchorRadius struct{
	Nid int
	AnchorRadius float64 //anchor的通信半径
}