--postgresql
-- Table: public.build_network

-- DROP TABLE public.build_network;

CREATE TABLE public.build_network
(
  nid integer NOT NULL DEFAULT nextval('build_network_bn_id_seq'::regclass), -- 主键
  bid integer NOT NULL, -- 楼宇id
  coor_id integer NOT NULL, -- 协调器id
  floor integer NOT NULL DEFAULT 1, -- 楼层
  status integer NOT NULL DEFAULT 1, -- 状态：0：关闭(维护中) 1：打开 2：异常
  CONSTRAINT build_network_pkey PRIMARY KEY (nid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.build_network
  OWNER TO "wsnDev";
COMMENT ON TABLE public.build_network
  IS '楼宇和网络表';
COMMENT ON COLUMN public.build_network.nid IS '主键';
COMMENT ON COLUMN public.build_network.bid IS '楼宇id';
COMMENT ON COLUMN public.build_network.coor_id IS '协调器id';
COMMENT ON COLUMN public.build_network.floor IS '楼层';
COMMENT ON COLUMN public.build_network.status IS '状态：0：关闭(维护中) 1：打开 2：异常';

-- Table: public.building

-- DROP TABLE public.building;

CREATE TABLE public.building
(
  bid integer NOT NULL DEFAULT nextval('building_bid_seq'::regclass), -- 建筑物id
  name character varying(30) NOT NULL DEFAULT ''::bpchar, -- 建筑物名称
  "position" point, -- 位置
  descrip character varying(100) NOT NULL DEFAULT ''::bpchar, -- 描述
  address character varying(50) NOT NULL DEFAULT ''::bpchar, -- 地址
  height integer NOT NULL DEFAULT 1, -- 高度
  CONSTRAINT building_pkey PRIMARY KEY (bid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.building
  OWNER TO "wsnDev";
COMMENT ON TABLE public.building
  IS '建筑物表';
COMMENT ON COLUMN public.building.bid IS '建筑物id';
COMMENT ON COLUMN public.building.name IS '建筑物名称';
COMMENT ON COLUMN public.building."position" IS '位置';
COMMENT ON COLUMN public.building.descrip IS '描述';
COMMENT ON COLUMN public.building.address IS '地址';
COMMENT ON COLUMN public.building.height IS '高度';



-- Table: public.building_map

-- DROP TABLE public.building_map;

CREATE TABLE public.building_map
(
  bm_id integer NOT NULL DEFAULT nextval('building_map_bm_id_seq'::regclass), -- 主键
  bid integer NOT NULL, -- 建筑id
  map_id integer NOT NULL, -- map id
  floor integer NOT NULL DEFAULT 1, -- 楼层
  CONSTRAINT building_map_pkey PRIMARY KEY (bm_id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.building_map
  OWNER TO "wsnDev";
COMMENT ON TABLE public.building_map
  IS 'building_map关系表';
COMMENT ON COLUMN public.building_map.bm_id IS '主键';
COMMENT ON COLUMN public.building_map.bid IS '建筑id';
COMMENT ON COLUMN public.building_map.map_id IS 'map id';
COMMENT ON COLUMN public.building_map.floor IS '楼层';

-- Table: public.map

-- DROP TABLE public.map;

CREATE TABLE public.map
(
  map_id integer NOT NULL DEFAULT nextval('map_map_id_seq'::regclass), -- 主键
  title character varying(40) NOT NULL DEFAULT ''::bpchar, -- 地图标题
  status integer NOT NULL DEFAULT 0, -- 状态，0：未发布 1：已发布 2：已锁定
  descrip character varying(50) DEFAULT ''::bpchar, -- 描述
  create_time timestamp without time zone NOT NULL DEFAULT now(),
  update_time timestamp without time zone NOT NULL DEFAULT now(),
  CONSTRAINT map_pkey PRIMARY KEY (map_id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.map
  OWNER TO "wsnDev";
COMMENT ON TABLE public.map
  IS '地图';
COMMENT ON COLUMN public.map.map_id IS '主键';
COMMENT ON COLUMN public.map.title IS '地图标题';
COMMENT ON COLUMN public.map.status IS '状态，0：未发布 1：已发布 2：已锁定';
COMMENT ON COLUMN public.map.descrip IS '描述';

-- Table: public.map_basemap

-- DROP TABLE public.map_basemap;

CREATE TABLE public.map_basemap
(
  map_id integer NOT NULL DEFAULT nextval('map_gis_map_id_seq'::regclass),
  code character varying(20) NOT NULL DEFAULT 'EPSG:404000'::character varying, -- 空间参考标识符,如SRID,ESPG。格式XXX:XX
  host character varying(100) NOT NULL DEFAULT ''::character varying, -- 地图服务所在的host
  server_type character varying(20) NOT NULL DEFAULT 'geoserver'::character varying, -- 地图服务类型
  workspace character varying(50) NOT NULL DEFAULT ''::character varying, -- 地图所在的服务工作空间
  request_type character varying(5) NOT NULL DEFAULT 'wms'::character varying, -- 地图请求类型，例如wms,wfs
  layers character varying(50) NOT NULL DEFAULT ''::character varying, -- 地图在服务中的图层名。
  x_min double precision NOT NULL DEFAULT 0, -- 地图范围x最小值
  x_max double precision NOT NULL DEFAULT 0,
  y_min double precision NOT NULL DEFAULT 0,
  y_max double precision NOT NULL DEFAULT 0,
  zoom_default double precision NOT NULL DEFAULT 1, -- 缩放，默认值
  zoom_max double precision NOT NULL DEFAULT 1,
  zoom_min double precision NOT NULL DEFAULT 1,
  CONSTRAINT map_gis_pkey PRIMARY KEY (map_id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.map_basemap
  OWNER TO "wsnDev";
COMMENT ON TABLE public.map_basemap
  IS '地图的底图信息表';
COMMENT ON COLUMN public.map_basemap.code IS '空间参考标识符,如SRID,ESPG。格式XXX:XX';
COMMENT ON COLUMN public.map_basemap.host IS '地图服务所在的host';
COMMENT ON COLUMN public.map_basemap.server_type IS '地图服务类型';
COMMENT ON COLUMN public.map_basemap.workspace IS '地图所在的服务工作空间';
COMMENT ON COLUMN public.map_basemap.request_type IS '地图请求类型，例如wms,wfs';
COMMENT ON COLUMN public.map_basemap.layers IS '地图在服务中的图层名。';
COMMENT ON COLUMN public.map_basemap.x_min IS '地图范围x最小值';
COMMENT ON COLUMN public.map_basemap.zoom_default IS '缩放，默认值';


-- Table: public.network_simu

-- DROP TABLE public.network_simu;

CREATE TABLE public.network_simu
(
  nid integer NOT NULL, -- 网络id
  anchor_radius double precision NOT NULL, -- anchor的通信半径
  CONSTRAINT network_simu_pkey PRIMARY KEY (nid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.network_simu
  OWNER TO "wsnDev";
COMMENT ON TABLE public.network_simu
  IS '用于模拟测试的网络表';
COMMENT ON COLUMN public.network_simu.nid IS '网络id';
COMMENT ON COLUMN public.network_simu.anchor_radius IS 'anchor的通信半径';



-- 插入测试数据
-- delete from public.building;
-- delete from public.building_map;
-- delete from public.map;
INSERT INTO public.building(
            bid, name, position, descrip, address,height)
    VALUES (1, '教学楼2号', point(114.271241,30.447683),'测试数据','湖北省武汉市某区某某大学',6);

INSERT INTO public.building_map(
            bm_id, bid, map_id, floor)
    VALUES (1, 1, 1, 4);

INSERT INTO public.map(
            map_id, title, status)
    VALUES (1, '教学楼2号-4F', 1);

INSERT INTO public.network(
            nid, sn, status)
    VALUES (1, '0x123456', 1);

INSERT INTO public.build_network(
            nid, bid, coor_id, floor, status)
    VALUES (1, 1, 1, 4, 1);

INSERT INTO public.network_simu(
            nid, anchor_radius)
    VALUES (1, 8);

INSERT INTO public.map_basemap(
            map_id, code, host, server_type, workspace, request_type, layers, 
            x_min, y_min,x_max, y_max, zoom_default, zoom_max, zoom_min)
    VALUES (1,'EPSG:404000','http://127.0.0.1:8083','geoserver','hzmap','wms'
    ,'gdata_1_1_plane',-10.153,4.394,77.846,71.617,1,5,1.2);


    