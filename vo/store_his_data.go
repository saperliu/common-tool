package vo

/**
 * 计算结果后的数据的存储的格式
 */
type StoreHisData struct {
	Tagid           string          `json:"tagid"`       //点位id
	Timestamp       string          `json:"timestamp"`   //记录时间
	Value           string          `json:"value"`       //点位值
	Status          int32           `json:"status"`      //状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界
	SiteId          string          `json:"site"`        //站点
	OrgId           string          `json:"org"`         //组织
	Tags            string          `json:"tags"`        //数据的标签
	StoreGPSHisData StoreGPSHisData `json:"gpsVo"`       //GPS 数据
	EnclosureVo     EnclosureVo     `json:"enclosureVo"` //围栏 数据
	DataType        string          `json:"dataType"`    //数据的类型
	DataDesc        string          `json:"dataDesc"`    //数据的描述
}

/**
 * 计算结果后的数据的存储的格式
 */
type StoreGPSHisData struct {
	Longitude    float64 `json:"longitude"`    //经度
	Latitude     float64 `json:"latitude"`     //纬度
	Speed        float64 `json:"speed"`        //速度 单位：km/h
	Direction    int32   `json:"direction"`    //方向 范围为[0,359]，0度为正北方向，顺时针
	Accuracy     float64 `json:"accuracy"`     //定位精度 单位：米
	Height       float64 `json:"height"`       //高度, 单位：米
	Status       int32   `json:"status"`       //状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界
	PositionType string  `json:"positionType"` //坐标类型：WIFI GPS  LBS等
}

//围栏数据的封装
type EnclosureVo struct {
	AssetId         string  `json:"assetId"`
	SerialNum       string  `json:"serialNum"`
	PersonId        string  `json:"personId"`
	CenterLongitude float64 `json:"centerLongitude"`
	CenterLatitude  float64 `json:"centerLatitude"`
	EffectiverAdius float64 `json:"effectiverAdius"`
}

/**
 * 计算结果后的数据的存储的格式
 */
type StoreCloudHisData struct {
	Tagid           string          `json:"tagid"`       //点位id
	Timestamp       float64         `json:"timestamp"`   //记录时间
	Value           string          `json:"value"`       //点位值
	Status          int32           `json:"status"`      //状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界
	SiteId          string          `json:"site"`        //站点
	OrgId           string          `json:"org"`         //组织
	Tags            string          `json:"tags"`        //数据的标签
	StoreGPSHisData StoreGPSHisData `json:"gpsVo"`       //GPS 数据
	EnclosureVo     EnclosureVo     `json:"enclosureVo"` //围栏 数据
	DataType        string          `json:"dataType"`    //数据的类型
	DataDesc        string          `json:"dataDesc"`    //数据的描述
}

/**
 * 自定义的带数据描述的数据
 */
type StorePerfectHisData struct {
	Value    string `json:"value"`    //点位值
	Status   int32  `json:"status"`   //状态 //0 正常,-1 离线, 1 通讯异常,2 数据越界
	DataType string `json:"dataType"` //数据的类型
	DataDesc string `json:"dataDesc"` //数据的描述
	Rvalue   string `json:"rvalue"`   //原始值
	Rstatus  int32  `json:"rstatus"`  //原始状态
}
