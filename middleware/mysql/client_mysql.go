package mysql

import (
	"common-tool/common"
	"common-tool/logger"
	"common-tool/vo"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

const (
	//user是AnalyticDB MySQL版集群中的用户账号：高权限账号或者普通账号。
	user = "adb_test"
	//password是AnalyticDB MySQL版集群中用户账号对应的密码。
	password = "xxx"
	//host是AnalyticDB MySQL版集群的连接地址，可以在控制台的集群信息页面获取连接地址。
	host = "127.0.xx.xx"
	//3306是端口号。
	port = 3306
	//database是AnalyticDB MySQL版集群中的数据库名称。
	database = "database_name"
	//数据库连接的超时时间。
	connectTimeout = "10s"
)

type Database struct {
	User         string
	Password     string
	Host         string
	Port         int
	Database     string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
	DB           *sql.DB
}

func (database *Database) CreateDB() (*sql.DB, error) {
	//打开数据库连接。
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&interpolateParams=true", database.User, database.Password, database.Host, database.Port, database.Database, connectTimeout)
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}
	//设置最大打开的连接数，默认值为0，表示不限制。
	db.SetMaxOpenConns(database.MaxOpenConns)
	//设置最大闲置的连接数。
	db.SetMaxIdleConns(database.MaxIdleConns)
	//设置连接的最大生命周期，默认是连接总是可重用。
	//该设置无法保证连接在连接池中完整存在一小时。连接可能会因为某些原因无法使用而自动关闭。
	//该设置项不是空闲超时时间，即连接会在第一次创建后一小时过期，而不是空闲后一小时过期。
	//理论上，连接的最大生命周期越短，从0开始创建连接的频率就会越高。
	db.SetConnMaxLifetime(time.Hour)
	//defer db.Close()
	//rows, err := db.Query("select * from student where id = ?", 9)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id string
	//	var name string
	//	var unit string
	//	err := rows.Scan(&id, &name, &unit)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	fmt.Println(fmt.Sprintf("%s, %s, %s", id, name, unit))
	//}
	database.DB = db
	return db, nil
}
func (database *Database) checkTableExist(table string) bool {
	rows, err := database.DB.Query("select count(*) from " + table + " ;")
	if err != nil {
		return false
	} else {
		defer rows.Close()
		return true
	}
}
func (database *Database) createTable(table string) bool {
	sql := fmt.Sprintf("CREATE TABLE `%s` (\n  `id` varchar(36) COLLATE utf8mb4_bin NOT NULL,\n  `timestamp` timestamp NULL DEFAULT NULL COMMENT '记录时间',\n  `value` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '点位值',\n  `longitude` float DEFAULT NULL COMMENT '经度',\n  `latitude` float DEFAULT NULL COMMENT '纬度',\n  `speed` float DEFAULT NULL COMMENT '速度 单位：km/h',\n  `direction` int DEFAULT NULL COMMENT '方向 范围为[0,359]，0度为正北方向，顺时针',\n  `accuracy` float DEFAULT NULL COMMENT '定位精度 单位：米',\n  `height` float DEFAULT NULL COMMENT '高度, 单位：米',\n  `status` int DEFAULT NULL COMMENT '状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界',\n  `site_id` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '站点',\n  `org_id` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组织',\n  `data_type` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '数据类型  normal  普通数据 enerbos GPS数据自定义数据  perfect 自定义的带数据描述的数据',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;", table)
	_, err := database.DB.Exec(sql)

	if err != nil {
		logger.Error("create  point table error: %s , table: %s", err.Error(), table)
		return false
	} else {
		return true
	}
}
func (database *Database) insertData(table string, val *vo.StorePoint, id string) bool {
	ti, err := strconv.ParseInt(val.Timestamp, 10, 64)
	if err != nil {
		logger.Error(" point.Timestamp  %v      %v ", val.Timestamp, err)
	}
	//时间是毫秒，转为纳秒再格式化时间
	dateTemp := time.Unix(0, ti*1e6)
	sql := fmt.Sprintf("INSERT INTO `%s` (`id`, `timestamp`, `value`, `longitude`, `latitude`, `speed`, `direction`, `accuracy`, `height`, `status`, `site_id`, `org_id`, `data_type`) VALUES ('%s', '%s', '%s', %f, %f, %f, %d, %f, %f, %d, '%s', '%s', '%s');", table, id, dateTemp.Format("2006-01-02 15:04:05"), val.Value, val.Longitude, val.Latitude, val.Speed, val.Direction, val.Accuracy, val.Height, val.Status, val.SiteId, val.OrgId, val.DataType)
	_, err1 := database.DB.Exec(sql)
	fmt.Printf(" insert  point  erro %s ", err1)
	if err1 != nil {
		logger.Error("insert  point = error: %s , table: %s,  data: %s", err.Error(), table, sql)
		return false
	} else {
		return true
	}
}

func (database *Database) insertHisData(table string, val *vo.StoreHisData, id string) bool {
	ti, err := strconv.ParseInt(val.Timestamp, 10, 64)
	if err != nil {
		logger.Error(" point.Timestamp  %v      %v ", val.Timestamp, err)
	}
	//时间是毫秒，转为纳秒再格式化时间
	dateTemp := time.Unix(0, ti*1e6)
	sql := fmt.Sprintf("INSERT INTO `%s` (`id`, `timestamp`, `value`, `longitude`, `latitude`, `speed`, `direction`, `accuracy`, `height`, `status`, `site_id`, `org_id`, `data_type`) VALUES ('%s', '%s', '%s', %f, %f, %f, %d, %f, %f, %d, '%s', '%s', '%s');", table, id, dateTemp.Format("2006-01-02 15:04:05"), val.Value, val.StoreGPSHisData.Longitude, val.StoreGPSHisData.Latitude, val.StoreGPSHisData.Speed, val.StoreGPSHisData.Direction, val.StoreGPSHisData.Accuracy, val.StoreGPSHisData.Height, val.Status, val.SiteId, val.OrgId, val.DataType)
	_, err1 := database.DB.Exec(sql)
	fmt.Printf(" insert  point  erro %s ", err1)
	if err1 != nil {
		logger.Error("insert  point = error: %s , table: %s,  data: %s", err.Error(), table, sql)
		return false
	} else {
		return true
	}
}

func (database *Database) SaveStorePoint(storePointVo vo.StorePoint) {

	if database.checkTableExist(database.TablePrefix + storePointVo.Tagid) {
		// 表存大直接保存数据
		database.insertData(database.TablePrefix+storePointVo.Tagid, &storePointVo, common.NumberToString(common.MustID()))
	} else {
		if database.createTable(database.TablePrefix + storePointVo.Tagid) {
			// 创建表成功，保存数据
			database.insertData(database.TablePrefix+storePointVo.Tagid, &storePointVo, common.NumberToString(common.MustID()))
		} else {
			// 创建表失败
			logger.Error("insert  point  error , table: %s,  data: %v ", storePointVo.Tagid, storePointVo)
		}
	}
}

func (database *Database) SaveHisPoint(storePointVo vo.StoreHisData) {

	if database.checkTableExist(database.TablePrefix + storePointVo.Tagid) {
		// 表存大直接保存数据
		database.insertHisData(database.TablePrefix+storePointVo.Tagid, &storePointVo, common.NumberToString(common.MustID()))
	} else {
		if database.createTable(database.TablePrefix + storePointVo.Tagid) {
			// 创建表成功，保存数据
			database.insertHisData(database.TablePrefix+storePointVo.Tagid, &storePointVo, common.NumberToString(common.MustID()))
		} else {
			// 创建表失败
			logger.Error("insert  point  error , table: %s,  data: %v ", storePointVo.Tagid, storePointVo)
		}
	}
}
