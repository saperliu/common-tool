package kairosdb_client

import (
	"github.com/json-iterator/go"
	"github.com/saperliu/common-tool/common"
	"github.com/saperliu/common-tool/logger"
	"github.com/saperliu/common-tool/vo"
	"github.com/saperliu/go-kairosdb/builder"
	"github.com/saperliu/go-kairosdb/client"
	"strconv"
	"strings"
	"time"
)

type KairosdbClient struct {
	Address  string //地址
	User     string //用户
	Password string //密码
	DbClient client.Client
}

var loggerKairosdb = logger.CreateLogger("kairosdb")

var WeekDayMap = map[string]string{
	"Monday":    "01",
	"Tuesday":   "02",
	"Wednesday": "03",
	"Thursday":  "04",
	"Friday":    "05",
	"Saturday":  "06",
	"Sunday":    "00",
}

func NewKairosdbClient(url string) *KairosdbClient {
	kairosdbClient := client.NewHttpClient(url)
	return &KairosdbClient{
		Address:  url,
		DbClient: kairosdbClient}
}

/**
 * 保存数据
 */
func (kairosDbClient *KairosdbClient) SaveHisPoint(storeHisData vo.StoreHisData) (error error) {

	//kairosdbClent := client.NewHttpClient(kairosDbClient.Address)
	mb := kairosDbClient.builderMetric(storeHisData)
	tempData, _ := jsoniter.Marshal(mb)
	loggerKairosdb.Info().Msgf(" SaveHisPoint  mb :  %s ", tempData)
	response, err := kairosDbClient.DbClient.PushMetrics(mb)
	//loggerKairosdb.Info().Msgf(" SaveHisPoint   %s    %v  if  err :  %v ", storeHisData.Tagid, response, err)
	if err != nil {
		return err
	} else {
		code := response.GetStatusCode()
		loggerKairosdb.Info().Msgf(" save point  code %v ", code)
	}
	return nil
}

func (kairosDbClient *KairosdbClient) builderMetric(storeHisData vo.StoreHisData) builder.MetricBuilder {

	//增加数据
	metricBuilder := builder.NewMetricBuilder()
	var metric builder.Metric
	timeTemp, err := strconv.ParseInt(storeHisData.Timestamp, 10, 64)
	if err != nil {
		logger.Error("----- storeHisData.Timestamp  %v      %v ", storeHisData.Timestamp, err)
	}
	dateTemp := time.Unix(0, timeTemp*1e6)

	if storeHisData.StoreGPSHisData.Latitude != 0 {
		// gps 数据
		metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storeHisData.Tagid)))
		metric.AddType(KAIROSDB_DATA_TYPE_GPS)
		if storeHisData.StoreGPSHisData.PositionType != "" {
			loggerKairosdb.Info().Msgf("----- PositionType  %v  ", storeHisData.StoreGPSHisData.PositionType)
			metric.AddTag(KAIROSDB_TAG_POSITION_TYPE, storeHisData.StoreGPSHisData.PositionType)
		}
		metric.AddDataPoint(timeTemp, storeHisData.StoreGPSHisData)
	} else {
		if storeHisData.DataType == KAIROSDB_DATA_TYPE_PERFECT {
			// 带有描述的数据
			metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storeHisData.Tagid)))
			metric.AddType(KAIROSDB_DATA_TYPE_PERFECT)
			if storeHisData.DataType != "" {
				loggerKairosdb.Info().Msgf("----- DataType  %v  ", storeHisData.DataType)
				metric.AddTag(KAIROSDB_TAG_DATA_TYPE, storeHisData.DataType)
			}
			var perfectData = vo.StorePerfectHisData{}
			perfectData.DataDesc = storeHisData.DataDesc
			perfectData.DataType = storeHisData.DataType
			perfectData.Value = storeHisData.Value
			perfectData.Status = storeHisData.Status
			perfectData.Rstatus = storeHisData.Status
			perfectData.Rvalue = storeHisData.Value

			metric.AddDataPoint(timeTemp, perfectData)

		} else {
			// 普通时序数据
			metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storeHisData.Tagid)))
			//boolNum :=util.IsNumeric(storeHisData.Value)
			metric.AddTag(KAIROSDB_TAG_RSTATUS, "0")
			v, err := strconv.ParseFloat(storeHisData.Value, 64)
			if err == nil {
				//v,err:=strconv.ParseFloat(storeHisData.Value,64)
				loggerKairosdb.Info().Msgf("----storeHisData  %s %s %v    ", storeHisData.Tagid, storeHisData.Value, v)
				metric.AddDataPoint(timeTemp, v)
			} else {
				loggerKairosdb.Info().Msgf("----storeHisData tagid %s    storeHisData.Value  %s    %v ", storeHisData.Tagid, storeHisData.Value, err)
				metric.AddDataPoint(timeTemp, storeHisData.Value)
			}
		}
	}

	if storeHisData.OrgId != "" {
		metric.AddTag(KAIROSDB_TAG_ORG, storeHisData.OrgId)
	}
	if storeHisData.SiteId != "" {
		metric.AddTag(KAIROSDB_TAG_SITE, storeHisData.SiteId)
	}
	metric.AddTag(KAIROSDB_TAG_YEAR, common.NumberToString(dateTemp.Year()))
	metric.AddTag(KAIROSDB_TAG_MONTH, dateTemp.Format("01"))
	metric.AddTag(KAIROSDB_TAG_DAY, dateTemp.Format("02"))
	metric.AddTag(KAIROSDB_TAG_HOUR, dateTemp.Format("15"))
	metric.AddTag(KAIROSDB_TAG_MIN, dateTemp.Format("04"))
	metric.AddTag(KAIROSDB_TAG_SEC, dateTemp.Format("05"))
	metric.AddTag(KAIROSDB_TAG_WEEK, WeekDayMap[dateTemp.Weekday().String()])
	if storeHisData.Tags != "" {
		metric.AddTag(KAIROSDB_TAG_REMARK, storeHisData.Tags)
	}
	return metricBuilder
}

/**
 * 保存数据
 */
func (kairosDbClient *KairosdbClient) SaveStorePoint(storePointVo vo.StorePoint) (error error) {

	//kairosdbClent := client.NewHttpClient(kairosDbClient.Address)
	mb := kairosDbClient.builderPointMetric(storePointVo)
	loggerKairosdb.Info().Msgf(" SaveStorePoint  storePointVo : %v ", storePointVo)
	response, err := kairosDbClient.DbClient.PushMetrics(mb)
	if err != nil {
		return err
	} else {
		code := response.GetStatusCode()
		loggerKairosdb.Info().Msgf(" SaveStorePoint  code %v ", code)
	}
	return nil
}

func (kairosDbClient *KairosdbClient) builderPointMetric(storePointVo vo.StorePoint) builder.MetricBuilder {

	//增加数据
	metricBuilder := builder.NewMetricBuilder()
	var metric builder.Metric
	timeTemp, err := strconv.ParseInt(storePointVo.Timestamp, 10, 64)
	if err != nil {
		logger.Error("builderPointMetric storeHisData.Timestamp  %v      %v ", storePointVo.Timestamp, err)
	}
	dateTemp := time.Unix(0, timeTemp*1e6)

	if storePointVo.DataType == KAIROSDB_DATA_TYPE_GPS {
		metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storePointVo.Tagid)))
		metric.AddType(KAIROSDB_DATA_TYPE_GPS)
		var gpsData = vo.StoreGPSHisData{}
		gpsData.Latitude = storePointVo.Latitude
		gpsData.Longitude = storePointVo.Longitude
		gpsData.Speed = storePointVo.Speed
		gpsData.Direction = storePointVo.Direction
		gpsData.Accuracy = storePointVo.Accuracy
		gpsData.Height = storePointVo.Height
		gpsData.Status = storePointVo.Status
		gpsData.PositionType = storePointVo.PositionType
		metric.AddDataPoint(timeTemp, gpsData)
	} else if storePointVo.DataType == KAIROSDB_DATA_TYPE_PERFECT {
		metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storePointVo.Tagid)))
		metric.AddType(KAIROSDB_DATA_TYPE_PERFECT)

		var perfectData = vo.StorePerfectHisData{}
		perfectData.DataDesc = storePointVo.DataDesc
		perfectData.DataType = storePointVo.DataType
		perfectData.Value = storePointVo.Value
		perfectData.Status = storePointVo.Status
		perfectData.Rstatus = storePointVo.Status
		perfectData.Rvalue = storePointVo.Value

		metric.AddDataPoint(timeTemp, perfectData)
	} else {
		metric = metricBuilder.AddMetric(strings.ToUpper(strings.TrimSpace(storePointVo.Tagid)))
		metric.AddTag(KAIROSDB_TAG_RSTATUS, "0")
		//boolNum :=util.IsNumeric(storePointVo.Value)
		v, err := strconv.ParseFloat(storePointVo.Value, 64)
		if err != nil {
			//v,err:=strconv.ParseFloat(storePointVo.Value,64)
			loggerKairosdb.Info().Msgf("----storePointVo %v    err  %v     ", v, err)
			metric.AddDataPoint(timeTemp, v)
		} else {
			metric.AddDataPoint(timeTemp, storePointVo.Value)
		}
	}

	if storePointVo.OrgId != "" {
		loggerKairosdb.Info().Msgf("----- OrgId  %v     ", storePointVo.OrgId)
		metric.AddTag(KAIROSDB_TAG_ORG, storePointVo.OrgId)
	}
	if storePointVo.SiteId != "" {
		loggerKairosdb.Info().Msgf("----- SiteId  %v  ", storePointVo.SiteId)
		metric.AddTag(KAIROSDB_TAG_SITE, storePointVo.SiteId)
	}
	if storePointVo.PositionType != "" {
		loggerKairosdb.Info().Msgf("----- PositionType  %v  ", storePointVo.PositionType)
		metric.AddTag(KAIROSDB_TAG_POSITION_TYPE, storePointVo.PositionType)
	}
	metric.AddTag(KAIROSDB_TAG_YEAR, common.NumberToString(dateTemp.Year()))
	metric.AddTag(KAIROSDB_TAG_MONTH, dateTemp.Format("01"))
	metric.AddTag(KAIROSDB_TAG_DAY, dateTemp.Format("02"))
	metric.AddTag(KAIROSDB_TAG_HOUR, dateTemp.Format("15"))
	metric.AddTag(KAIROSDB_TAG_MIN, dateTemp.Format("04"))
	metric.AddTag(KAIROSDB_TAG_SEC, dateTemp.Format("05"))
	metric.AddTag(KAIROSDB_TAG_WEEK, WeekDayMap[dateTemp.Weekday().String()])
	return metricBuilder
}

func (kairosDbClient *KairosdbClient) QueryMetric(metricName string, startDate int64, endDate int64) []vo.StorePoint {

	//增加数据
	metricBuilder := builder.NewQueryBuilder()
	metricBuilder.AddMetric(metricName)
	//gpsTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-01-28 19:01:00", time.Local)
	dateTemp := time.Unix(0, startDate*1e6)
	dateTemp2 := time.Unix(0, endDate*1e6)
	metricBuilder.SetAbsoluteStart(dateTemp.Local())
	metricBuilder.SetAbsoluteEnd(dateTemp2.Local())
	metricBuilder.SetCacheTime(0)

	loggerKairosdb.Info().Msgf("------------  %v ", metricBuilder)
	//cli := client.NewHttpClient(kairosDbClient.Address)
	jsonUtil := jsoniter.ConfigCompatibleWithStandardLibrary
	queryResp, _ := kairosDbClient.DbClient.Query(metricBuilder)
	var listData []vo.StorePoint
	for _, tagName := range queryResp.QueriesArr {
		loggerKairosdb.Info().Msgf("--------tagName----  %v ", tagName)
		for _, tagPoint := range tagName.ResultsArr {

			loggerKairosdb.Info().Msgf("--------tagPoint----  %v ", tagPoint)
			for _, point := range tagPoint.DataPoints {

				gpsData := vo.StorePoint{}
				gpsData.Tagid = metricName
				gpsData.Timestamp = common.NumberToString(point.Timestamp())
				tempByte, _ := point.ByteValue()
				err := jsonUtil.Unmarshal(tempByte, &gpsData)
				if err != nil {
					logger.Error("--------point----  %v ", gpsData)
				}
				loggerKairosdb.Info().Msgf("--------point----  %v ", point)
				listData = append(listData, gpsData)
			}
		}
	}
	loggerKairosdb.Info().Msgf("--------listData----  %v ", listData)
	return listData
}
