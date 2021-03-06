package sensor

// 格式：每五个字段换行。
//当被连表数据不存在时，采用left join才能出现结果
const (
	getXYOfAnchorSqlTemp = `
	SELECT 
		gid,layer as type ,st_x(geom)as x,st_y(geom)as y 
	FROM 
		gsimu_%d_anchor
	WHERE 
		layer=$1
	`
)

const (
	getAnchorRadiusSql = `
	SELECT 
		nid,anchor_radius
	FROM 
		network_simu
	WHERE 
		nid=$1
	`
)

const(
	updateNetworkStatusSql=`
	UPDATE
		build_network
	SET
		status=$1
	WHERE
		nid=$2 and floor=$3
	`
)