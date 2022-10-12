package ali_tsdb

type AliMetric struct {
	Metric    string            `json:"metric"`    //测点
	Timestamp int64             `json:"timestamp"` //记录时间
	Value     float64           `json:"value"`     //点位值
	Tags      map[string]string `json:"tags"`      //测点的标签
}
