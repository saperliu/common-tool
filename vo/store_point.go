package vo

/**
 * 计算结果后的数据的存储的格式
 */
type StorePoint struct {
	Id              string  `json:"id"`              //数据的id，mysql存储时使用
	Tagid           string  `json:"tagid"`           //点位id
	Timestamp       string  `json:"timestamp"`       //记录时间
	Value           string  `json:"value"`           //点位值
	Longitude       float64 `json:"longitude"`       //经度
	Latitude        float64 `json:"latitude"`        //纬度
	Speed           float64 `json:"speed"`           //速度 单位：km/h
	Direction       int32   `json:"direction"`       //方向 范围为[0,359]，0度为正北方向，顺时针
	Accuracy        float64 `json:"accuracy"`        //定位精度 单位：米
	Height          float64 `json:"height"`          //高度, 单位：米
	Status          int32   `json:"status"`          //状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界
	SiteId          string  `json:"siteId"`          //站点
	OrgId           string  `json:"orgId"`           //组织
	DataType        string  `json:"dataType"`        //数据类型  normal  普通数据 gps GPS数据自定义数据  perfect 自定义的带数据描述的数据
	PositionType    string  `json:"positionType"`    //坐标类型：WIFI GPS  LBS等
	DataPerfectType string  `json:"dataPerfectType"` //数据的类型,默认:DEFAULT
	DataDesc        string  `json:"dataDesc"`        //数据的描述
}
