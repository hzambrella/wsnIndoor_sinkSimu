--postgresql

-- Table: public.building

-- DROP TABLE public.building;

CREATE TABLE public.building
(
  bid integer NOT NULL DEFAULT nextval('building_bid_seq'::regclass), -- 建筑物id
  name character(30) NOT NULL DEFAULT ''::bpchar, -- 建筑物名称
  "position" point, -- 位置
  descrip character(1) NOT NULL DEFAULT ''::bpchar, -- 描述
  address character(50) NOT NULL DEFAULT ''::bpchar, -- 地址
  CONSTRAINT building_pkey PRIMARY KEY (bid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.building
  OWNER TO postgres;
COMMENT ON TABLE public.building
  IS '建筑物表';
COMMENT ON COLUMN public.building.bid IS '建筑物id';
COMMENT ON COLUMN public.building.name IS '建筑物名称';
COMMENT ON COLUMN public.building."position" IS '位置';
COMMENT ON COLUMN public.building.descrip IS '描述';
COMMENT ON COLUMN public.building.address IS '地址';


-- Table: public.building_map

-- DROP TABLE public.building_map;

CREATE TABLE public.building_map
(
  bm_id serial NOT NULL, -- 主键
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
  map_id serial NOT NULL, -- 主键
  title character(40) NOT NULL DEFAULT ''::bpchar, -- 地图标题
  status integer NOT NULL DEFAULT 0, -- 状态，0：未发布 1：已发布 2：已锁定
  descrip character(50) DEFAULT ''::bpchar, -- 描述
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


-- Table: public.network

-- DROP TABLE public.network;

CREATE TABLE public.network
(
  nid integer NOT NULL DEFAULT nextval('network_nid_seq'::regclass), -- 主键，网络id
  sn character(18) NOT NULL DEFAULT ''::bpchar, -- 网关设备序列号
  status integer NOT NULL DEFAULT 1, -- 状态 0-关闭 1-打开 2-维护中
  CONSTRAINT network_pkey PRIMARY KEY (nid)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.build_network
  OWNER TO "wsnDev";
COMMENT ON TABLE public.network
  IS '网络信息表';
COMMENT ON COLUMN public.network.nid IS '主键，网络id';
COMMENT ON COLUMN public.network.sn IS '网关设备序列号';
COMMENT ON COLUMN public.network.status IS '状态 0-关闭 1-打开 2-维护中';

-- Table: public.build_network

-- DROP TABLE public.build_network;

CREATE TABLE public.build_network
(
  bn_id integer NOT NULL DEFAULT nextval('build_network_bn_id_seq'::regclass), -- 主键
  bid integer NOT NULL, -- 楼宇id
  coor_id integer NOT NULL, -- 协调器id
  floor integer NOT NULL DEFAULT 1, -- 楼层
  status integer NOT NULL DEFAULT 1, -- 状态：0：关闭(维护中) 1：打开 2：异常
  CONSTRAINT build_network_pkey PRIMARY KEY (bn_id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.build_network
  OWNER TO "wsnDev";
COMMENT ON TABLE public.build_network
  IS '楼宇和网络表';
COMMENT ON COLUMN public.build_network.bn_id IS '主键';
COMMENT ON COLUMN public.build_network.bid IS '楼宇id';
COMMENT ON COLUMN public.build_network.coor_id IS '协调器id';
COMMENT ON COLUMN public.build_network.floor IS '楼层';
COMMENT ON COLUMN public.build_network.status IS '状态：0：关闭(维护中) 1：打开 2：异常';

-- Table: public.network_simu

-- DROP TABLE public.network_simu;

CREATE TABLE public.network_simu
(
  nid integer NOT NULL, -- 网络id
  anchor_radius double precision NOT NULL, -- anchor的通信半径
  floor integer NOT NULL DEFAULT 1, -- 楼层
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
COMMENT ON COLUMN public.network_simu.floor IS '楼层';


-- 插入测试数据
-- delete from public.building;
-- delete from public.building_map;
-- delete from public.map;
INSERT INTO public.building(
            bid, name, position, descrip, address)
    VALUES (1, '教学楼2号', point(114.271241,30.447683),'','');

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
            bn_id, bid, coor_id, floor, status)
    VALUES (1, 1, 1, 4, 1);

INSERT INTO public.network_simu(
            nid, anchor_radius, floor)
    VALUES (1, 8, 4);




    